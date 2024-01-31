<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;
    private array $data;

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $this->data = (array)$this->getRequest('data');
        $result = $this->notification_structure();

        if (!$this->validation_request()) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $result;
        }

        $members = $this->getMembers();
        $notificationType = (int)$this->data['notificationType'];

        $templateData = $this->teamplate_data($members);
        $this->validation_teamplate_params($templateData);

        return $this->sendEmails($members, $templateData, $result, $notificationType);
    }

    private function __(string $string, array $templateData, $resellerId)
    {

    }

    /**
     * @return array
     */
    public function notification_structure(): array
    {
        return [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];
    }

    private function validation_request(): bool
    {
        if (empty($this->data['resellerId'])) {
            return false;
        }

        $required_params = [
            'notificationType' => 'Empty notificationType',
            'resellerId' => 'Seller not found!',
            'clientId' => 'Client not found!',
            'creatorId' => 'Creator not found!',
            'expertId' => 'Expert not found!',
        ];

        foreach ($required_params as $field => $message) {
            if (empty($data[$field])) {
                throw new \Exception($message, 400);
            }
        }

        return true;
    }

    /**
     * @return Members
     * @throws \Exception
     */
    private function getMembers(): Members
    {
        $reseller = Seller::getById((int)$this->data['resellerId']);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        $client = Contractor::getById((int)$this->data['clientId']);
        if ($client === null || $client->getType() !== Contractor::TYPE_CUSTOMER || $client->Seller->getId() !== $reseller->Seller->getId()) { //todo предположим что у объекта есть Seller
            throw new \Exception('сlient not found!', 400);
        }

        $creator = Employee::getById((int)$this->data['creatorId']);
        if ($creator === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $expert = Employee::getById((int)$this->data['expertId']);
        if ($expert === null) {
            throw new \Exception('Expert not found!', 400);
        }

        return new Members($reseller, $client, $creator, $expert);
    }

    /**
     * @param int $notificationType
     * @param int $reseller_id
     * @return array
     */
    public function differences(int $notificationType, int $reseller_id): array
    {
        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $reseller_id);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)($data['differences']['from'] ?? 0)),
                'TO' => Status::getName((int)($data['differences']['to'] ?? 0)),
            ], $reseller_id);
        }
        return $differences;
    }

    /**
     * @param Members $members
     * @return array
     */
    private function teamplate_data(Members $members): array
    {
        $notificationType = (int)$this->data['notificationType'] ?? 0;

        $differences = $this->differences($notificationType, $members->reseller->getId());
        $cFullName = $members->client->getFullName();
        if (empty($members->client->getFullName())) {
            $cFullName = $members->client->getName();
        }

        return [
            'COMPLAINT_ID' => (int)$this->data['complaintId'] ?? '',
            'COMPLAINT_NUMBER' => (string)$this->data['$complaintNumber'] ?? '',
            'CREATOR_ID' => $members->creator->getId(),
            'CREATOR_NAME' => $members->creator->getFullName(),
            'EXPERT_ID' => $members->expert->getId(),
            'EXPERT_NAME' => $members->expert->getFullName(),
            'CLIENT_ID' => $members->expert->getId(),
            'CLIENT_NAME' => $cFullName,
            'CONSUMPTION_ID' => (int)$this->data['consumptionId'] ?? '',
            'CONSUMPTION_NUMBER' => (string)$this->data['consumptionNumber'] ?? '',
            'AGREEMENT_NUMBER' => (string)$this->data['agreementNumber'] ?? '',
            'DATE' => (string)$this->data['date'] ?? '',
            'DIFFERENCES' => $differences,
        ];
    }

    /**
     * @param array $templateData
     * @return void
     * @throws \Exception
     */
    public function validation_teamplate_params(array $templateData): void
    {
        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception('Template Data (' . $key . ') is empty!', 500);
            }
        }
    }

    /**
     * @param string $emailFrom
     * @param string $email
     * @param array $templateData
     * @param int $resellerId
     * @return bool
     */
    public function SendEmployeeEmail(string $emailFrom, string $email, array $templateData, int $resellerId): bool
    {
        MessagesClient::sendMessage([
            0 => [ // MessageTypes::EMAIL
                'emailFrom' => $emailFrom,
                'emailTo' => $email,
                'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
            ],
        ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);

        return true;
    }

    /**
     * @param string $emailFrom
     * @param Contractor $client
     * @param array $templateData
     * @param int $resellerId
     * @param int $differences
     * @return bool
     */
    public function SendClientEmail(string $emailFrom, Contractor $client, array $templateData, int $resellerId, int $differences): bool
    {
        MessagesClient::sendMessage([
            0 => [ // MessageTypes::EMAIL
                'emailFrom' => $emailFrom,
                'emailTo' => $client->getEmail(),
                'subject' => $this->__('complaintClientEmailSubject', $templateData, $resellerId),
                'message' => $this->__('complaintClientEmailBody', $templateData, $resellerId),
            ],
        ], $resellerId, $client->getID(), NotificationEvents::CHANGE_RETURN_STATUS, (int)$differences);

        return true;

    }

    /**
     * @param Members $members
     * @param array $templateData
     * @param array $result
     * @param int $notificationType
     * @return array
     */
    private function sendEmails(Members $members, array $templateData, array $result, int $notificationType): array
    {
        $emailFrom = getResellerEmailFrom();
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($members->reseller->getId(), 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                $result['notificationEmployeeByEmail'] = $this->SendEmployeeEmail($emailFrom, $email, $templateData, $members->reseller->getId());
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($members->client->getEmail())) {
                $result['notificationClientByEmail'] = $this->SendClientEmail($emailFrom, $members->client, $templateData, $members->reseller->getId(), (int)$data['differences']['to']);
            }

            if (!empty($members->client->getMobile())) {
                $result = $this->SendNotificationManager($members, (int)$data['differences']['to'], $templateData, $result);
            }
        }

        return $result;
    }

    /**
     * @param Members $members
     * @param int $differences
     * @param array $templateData
     * @param array $result
     * @return array
     */
    public function SendNotificationManager(Members $members, int $differences, array $templateData, array $result): array
    {
        $error = '';
        $res = NotificationManager::send($members->reseller->getId(), $members->client->getId(), NotificationEvents::CHANGE_RETURN_STATUS, $differences, $templateData, $error);  //TODO допустим, что $error возвращается по ссылке
        if ($res) {
            $result['notificationClientBySms']['isSent'] = true;
        }
        if (!empty($error)) {
            $result['notificationClientBySms']['message'] = $error;
        }

        return $result;
    }
}

class Members
{
    /**
     * @var Contractor
     */
    public Contractor $reseller;
    public Contractor $client;
    public Contractor $creator;
    public Contractor $expert;

    /**
     * @param Contractor $reseller
     * @param Contractor $client
     * @param Contractor $creator
     * @param Contractor $expert
     */
    public function __construct(Contractor $reseller, Contractor $client, Contractor $creator, Contractor $expert)
    {
        $this->reseller = $reseller;
        $this->client = $client;
        $this->creator = $creator;
        $this->expert = $expert;
    }
}

