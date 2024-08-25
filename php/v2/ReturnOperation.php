<?php

namespace NW\WebService\References\Operations\Notification;

class ReturnOperation extends ReferencesOperation
{
    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data = (array) $this->getRequest('data');

        $this->validateData($data);

        $notification = $this->mapArrayToNotification($data);

        $result = $this->initializeResult();

        $this->ensureSellerExistOrFail($notification->getResellerId());

        $client = $this->findContractor($data['clientId']);
        $this->validateClient($client, $notification);

        $creatorEmployee = $this->findEmployee((int) $data['creatorId']);
        $expertEmployee = $this->findEmployee((int) $data['expertId']);

        $differences = $this->getDifferences($notification);

        $templateData = new NotificationTemplateData();
        $templateData->fillWithNotification($notification);
        $templateData->fillWithCreator($creatorEmployee);
        $templateData->fillWithExpert($expertEmployee);
        $templateData->fillWithClient($client);
        $templateData->fillWithDifferences($differences);

        $templateData->validate();

        $templateData = $templateData->getTemplateData();

        $emailFrom = getResellerEmailFrom($notification->getResellerId());
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($notification->getResellerId(), 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            $this->trySendEmployeeEmails($emailFrom, $emails, $templateData, $notification->getResellerId());
            $result['notificationEmployeeByEmail'] = true;
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notification->getType() === NotificationDTO::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                $this->trySendClientNotification($emailFrom, $client, $templateData, $notification->getResellerId());

                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                try {
                    $this->trySendMobileNotification(
                        $notification->getResellerId(),
                        $client->getId(),
                        (int) $data['differences']['to'],
                        $templateData
                    );

                    $result['notificationClientBySms']['isSent'] = true;
                } catch (MessageNotSentException $exception) {
                    $result['notificationClientBySms']['isSent'] = false;
                } catch (\Exception $exception) {
                    $result['notificationClientBySms']['message'] = $exception->getMessage();
                }
            }
        }

        return $result;
    }

    private function mapArrayToNotification(array $data): NotificationDTO
    {
        $notification = new NotificationDTO();
        $notification
            ->setResellerId((int) $data['resellerId'])
            ->setType((int) $data['notificationType'])
            ->setClientId((int) $data['clientId'])
            ->setCreatorId((int) $data['creatorId'])
            ->setExpertId((int) $data['expertId'])
            ->setComplaintId((int) $data['complaintId'])
            ->setComplaintNumber((string) $data['complaintNumber'])
            ->setConsumptionId((int) $data['consumptionId'])
            ->setConsumptionNumber((string) $data['consumptionNumber'])
            ->setAgreementNumber((string) $data['agreementNumber'])
            ->setDate((string) $data['date'])
            ->setDifferences($data['differences'] ?? []);

        return $notification;
    }

    private function initializeResult(): array
    {
        return [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];
    }

    private function validateData(array $data): void
    {
        if (empty($data['resellerId'])) {
            throw new \Exception('Empty resellerId', Response::HTTP_BAD_REQUEST);
        }

        if (empty($data['notificationType'])) {
            throw new \Exception('Empty notificationType', Response::HTTP_BAD_REQUEST);
        }
    }

    private function ensureSellerExistOrFail(int $resellerId): void
    {
        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found! id: ' . $resellerId, Response::HTTP_NOT_FOUND);
        }
    }

    private function validateClient(Contractor $client, NotificationDTO $notificationDTO)
    {
        if (!$client->isCustomer() || $client->getSeller()->getId() !== $notificationDTO->getResellerId()) {
            throw new HttpInternalServerErrorException('Client not found! id: ' . $client->getId());
        }
    }

    private function findContractor(int $contractorId): Contractor
    {
        $contractor = Contractor::getById($contractorId);
        if ($contractor === null) {
            throw new \Exception('Contractor not found! id: ' . $contractorId, Response::HTTP_NOT_FOUND);
        }

        return $contractor;
    }

    private function findEmployee(int $employeeId): Employee
    {
        $employee = Employee::getById($employeeId);
        if ($employee === null) {
            throw new \Exception('Employee not found!', Response::HTTP_NOT_FOUND);
        }

        return $employee;
    }

    private function getDifferences(NotificationDTO $notification)
    {
        switch ($notification->getType()) {
            case NotificationDTO::TYPE_NEW:
                return __('NewPositionAdded', null, $notification->getResellerId());
            case NotificationDTO::TYPE_CHANGE:
                return __('PositionStatusHasChanged', [
                    'FROM' => Status::getStatus((int) $notification->getDifferences()['from']),
                    'TO'   => Status::getStatus((int) $notification->getDifferences()['to']),
                ], $notification->getResellerId());
            default:
                throw new HttpBadRequestException('Unsupported notification type: ' . $notification->getType());
        }
    }

    private function trySendEmployeeEmails(string $emailFrom, array $emails, array $templateData, int $resellerId): void
    {
        foreach ($emails as $email) {
            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $email,
                    'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
        }
    }

    private function trySendClientNotification($emailFrom, Contractor $client, $templateData, $resellerId)
    {
        MessagesClient::sendMessage([
            0 => [ // MessageTypes::EMAIL
                'emailFrom' => $emailFrom,
                'emailTo'   => $client->email,
                'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
            ],
        ],
            $resellerId,
            $client->getId(),
            NotificationEvents::CHANGE_RETURN_STATUS,
            $templateData['DIFFERENCES']['to']
        );
    }

    private function trySendMobileNotification(int $resellerId, int $clientId, int $status, array $templateData)
    {
        $result = NotificationManager::send(
            $resellerId,
            $clientId,
            NotificationEvents::CHANGE_RETURN_STATUS,
            $status,
            $templateData,
            $error
        );

        if (!$result) {
            throw new MessageNotSentException();
        }

        if (!empty($error)) {
            throw new \Exception($error);
        }
    }
}
