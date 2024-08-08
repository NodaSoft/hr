<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    private const TYPE_NEW    = 1;
    private const TYPE_CHANGE = 2;

    /**
     * @throws \Exception - не удалось создать объект Contractor
     * @throws \Exception - переменная для шаблона не задана
     */
    public function doOperation(): array
    {
        $data = $this->getRequest('data');
        $result = [
          'notificationEmployeeByEmail' => false,
          'notificationClientByEmail'   => false,
          'notificationClientBySms'     => [
            'isSent'  => false,
            'error' => '',
          ],
        ];

        if (empty($data)) {
            return $result;
        }

        $reseller = $this->createContractor(Contractor::TYPE_SELLER, $data['resellerId']);

        if (! $reseller) {
            $result['notificationClientBySms']['error'] = 'Empty resellerId';
            return $result;
        }

        $notificationType = filter_var($data['notificationType'], FILTER_VALIDATE_INT);
        if (! $notificationType) {
            throw new \Exception('Empty notificationType', 400);
        }

        $client = $this->createContractor(Contractor::TYPE_CLIENT, $data['clientId']);
        if (! $client) {
            throw new \Exception('сlient not found!', 400);
        }

        if ($client->getSeller()->getId() !== $reseller->getId()) {
            throw new \Exception('Invalid seller for client', 400);
        }

        $creator = $this->createContractor(Contractor::TYPE_EMPLOYEE, $data['creatorId']);
        if (! $creator) {
            throw new \Exception('Creator not found!', 400);
        }

        $expert = $this->createContractor(Contractor::TYPE_EMPLOYEE, $data['expertId']);
        if (! $expert) {
            throw new \Exception('Expert not found!', 400);
        }

        $differences = '';
        $notificationTo = false;
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $reseller->getId());
        } elseif ($notificationType === self::TYPE_CHANGE) {
            $notificationFrom = filter_var($data['differences']['from'], FILTER_VALIDATE_INT);
            $notificationTo = filter_var($data['differences']['to'], FILTER_VALIDATE_INT);

            if ($notificationTo && $notificationFrom) {
                $differences = __('PositionStatusHasChanged', [
                  'FROM' => Status::getName($notificationFrom),
                  'TO'   => Status::getName($notificationTo),
                ], $reseller->getId());
            }
        }

        $templateData = [
          'COMPLAINT_ID'       => filter_var($data['complaintId'], FILTER_VALIDATE_INT),
          'COMPLAINT_NUMBER'   => filter_var($data['complaintNumber']),
          'CREATOR_ID'         => $creator->getId(),
          'CREATOR_NAME'       => trim($creator->getFullName()),
          'EXPERT_ID'          => $expert->getId(),
          'EXPERT_NAME'        => trim($expert->getFullName()),
          'CLIENT_ID'          => $client->getId(),
          'CLIENT_NAME'        => trim($client->getFullName()),
          'CONSUMPTION_ID'     => filter_var($data['consumptionId'], FILTER_VALIDATE_INT),
          'CONSUMPTION_NUMBER' => filter_var($data['consumptionNumber']),
          'AGREEMENT_NUMBER'   => filter_var($data['agreementNumber']),
          'DATE'               => filter_var($data['date']),
          'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $templateKey => $templateItem) {
            if (empty($templateItem)) {
                throw new \Exception("Template Data ({$templateKey}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom();
        if (empty($emailFrom)) {
            return $result;
        }
        // Получаем email сотрудников из настроек
        $employeeEmails = getEmailsByPermit($reseller->getId(), 'tsGoodsReturn');
        if (! empty($employeeEmails)) {
            $employeeMessageSent = false;
            foreach ($employeeEmails as $employeeEmail) {
                $employeeMessageSent = $employeeMessageSent || MessagesClient::sendMessage([
                    [ // MessageTypes::EMAIL
                      'emailFrom' => $emailFrom,
                      'emailTo'   => $employeeEmail,
                      'subject'   => __('complaintEmployeeEmailSubject', $templateData, $reseller->getId()),
                      'message'   => __('complaintEmployeeEmailBody', $templateData, $reseller->getId()),
                    ],
                  ], $reseller->getId(), NotificationEvents::NEW_RETURN_STATUS->value);
            }

            if ($employeeMessageSent) {
                $result['notificationEmployeeByEmail'] = true;
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && $notificationTo) {
            if (! empty($client->getEmail())) {
                $clientMessageSent = MessagesClient::sendMessage([
                  [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $client->getEmail(),
                    'subject'   => __('complaintClientEmailSubject', $templateData, $reseller->getId()),
                    'message'   => __('complaintClientEmailBody', $templateData, $reseller->getId()),
                  ],
                ], $reseller->getId(), $client->getId(), NotificationEvents::CHANGE_RETURN_STATUS->value, $notificationTo);

                if ($clientMessageSent) {
                    $result['notificationClientByEmail'] = true;
                }
            }

            if (! empty($client->getMobile())) {
                try {
                    $notificationSent = NotificationManager::send($reseller->getId(), $client->getId(), NotificationEvents::CHANGE_RETURN_STATUS->value,
                      $notificationTo, $templateData);
                    if ($notificationSent) {
                        $result['notificationClientBySms']['isSent'] = true;
                    }
                } catch (\Exception $exception) {
                    $result['notificationClientBySms']['error'] = $exception->getMessage();
                }
            }
        }

        return $result;
    }

    private function createContractor(int $type, $id): Contractor|false
    {
        $id = filter_var($id, FILTER_VALIDATE_INT);
        $contractor = Contractor::getById($id);

        if ($contractor && $contractor->getType() === $type) {
            return $contractor;
        }

        return false;
    }
}
