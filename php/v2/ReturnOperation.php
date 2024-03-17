<?php

namespace NW\WebService\References\Operations\Notification\Operation;


class TsReturnOperation extends ReferencesOperation
{

    private array $data;
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    public function __construct()
    {
        $this->data = $this->getRequest('data');
    }

    public function doOperation(): array
    {
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        $resellerId = $this->getResselerId();
        $notificationType = $this->getNotificationType();

        $this->validateData();

        $reseller = $this->findReseller($resellerId);
        $client = $this->findClient($resellerId);

        $cFullName = $client->getFullName();
        if (empty($cFullName)) {
            $cFullName = $client->name;
        }

        $cr = $this->findEmployee($data['creatorId']);
        $et = $this->findEmployee($data['expertId']);

        $differences = $this->calculateDifferences($notificationType);

        $templateData = $this->buildTemplateData($data, $cr, $et, $cFullName, $differences);

        $this->validateTemplateData($templateData);

        $emailFrom = $this->getResellerEmailFrom($resellerId);
        $this->sendEmployeeNotifications($emailFrom, $templateData, $resellerId);

        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            $this->sendClientNotifications($emailFrom, $templateData, $resellerId, $client, $data);
        }

        return $result;
    }

    private function getResselerId(): int
    {
        if (!isset($this->data['resellerId'])) {
            throw new \Exception('Empty resellerId', 400);
        }
        return $this->data['resellerId'];
    }

    private function getNotificationType(): int
    {
        if (!isset($this->data['notificationType'])) {
            throw new \Exception('Empty notificationType', 400);
        }
        return $this->data['notificationType'];
    }

    private function validateData(): void
    {
        $resellerId = $this->getResselerId();
        if (empty($resellerId)) {
            throw new \Exception('Empty resellerId', 400);
        }

        $notificationType = $this->getNotificationType();
        if (empty($notificationType)) {
            throw new \Exception('Empty notificationType', 400);
        }
    }

    private function findReseller(int $resellerId)
    {
        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }
        return $reseller;
    }

    private function findClient(int $resellerId): Client
    {
        $client = Contractor::getById($this->data['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Client not found or invalid!', 400);
        }
        return $client;
    }

    private function findEmployee(int $employeeId)
    {
        $employee = Employee::getById($employeeId);
        if ($employee === null) {
            throw new \Exception('Employee not found!', 400);
        }
        return $employee;
    }

    private function calculateDifferences(int $notificationType): string
    {
        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($this->data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$this->data['differences']['from']),
                'TO'   => Status::getName((int)$this->data['differences']['to']),
            ], $resellerId);
        }
        return $differences;
    }

    private function buildTemplateData(array $data, $cr, $et, string $cFullName, string $differences): array
    {
        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => (int)$data['creatorId'],
            'CREATOR_NAME'       => $cr->getFullName(),
            'EXPERT_ID'          => (int)$data['expertId'],
            'EXPERT_NAME'        => $et->getFullName(),
            'CLIENT_ID'          => (int)$data['clientId'],
            'CLIENT_NAME'        => $cFullName,
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $differences,
        ];
        return $templateData;
    }

    private function validateTemplateData(array $templateData): void
    {
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }
    }

    private function getResellerEmailFrom(int $resellerId): string
    {
        $emailFrom = getResellerEmailFrom($resellerId);
        return $emailFrom;
    }

    private function sendEmployeeNotifications(string $emailFrom, array $templateData, int $resellerId): void
    {
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $email,
                        'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
            }
        }
    }

    private function sendClientNotifications(string $emailFrom, array $templateData, int $resellerId, Client $client, array $data): void
    {
        if (!empty($emailFrom) && !empty($client->email)) {
            MessagesClient::sendMessage([
                0 => [
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $client->email,
                    'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                    'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);
        }

        if (!empty($client->mobile)) {
            $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
            if ($res) {
                $result['notificationClientBySms']['isSent'] = true;
            }
            if (!empty($error)) {
                $result['notificationClientBySms']['message'] = $error;
            }
        }
    }
}
