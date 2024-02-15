<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\MessagesClient;
use NW\WebService\NotificationManager;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    private array $result;

    public function __construct()
    {
        $this->result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];
    }

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data = (array) $this->getRequest('data');
        $resellerId = $data['resellerId'];
        $notificationType = (int) $data['notificationType'];

        $this->validateInput($resellerId, $notificationType);

        $reseller = Seller::getById((int) $resellerId);
        $client = Contractor::getById((int) $data['clientId']);

        $this->validateEntities($reseller, $client, $resellerId);

        $templateData = $this->generateTemplateData($data, $resellerId, $client);

        $this->validateTemplateData($templateData);

        $emailFrom = getResellerEmailFrom($resellerId);

        $this->sendEmployeeNotifications($emailFrom, $templateData, $resellerId);
        $this->sendClientNotifications($notificationType, $data, $emailFrom, $client, $templateData, $resellerId);

        return $this->result;
    }

    private function validateInput( $resellerId, $notificationType): void
    {
        if (empty($resellerId) || empty($notificationType)) {
            throw new \Exception('Invalid input parameters', 400);
        }
    }

    private function validateEntities($reseller, $client, $resellerId): void
    {
        if ($reseller === null || $client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Invalid entities', 400);
        }
    }

    private function generateTemplateData(array $data, int $resellerId, Contractor $client): array
    {
        $differences = '';
        $creator = Employee::getById((int) $data['creatorId']);
        $expert = Employee::getById((int) $data['expertId']);

        if ($data['notificationType'] === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($data['notificationType'] === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int) $data['differences']['from']),
                'TO'   => Status::getName((int) $data['differences']['to']),
            ], $resellerId);
        }

        return [
            'COMPLAINT_ID'       => (int) $data['complaintId'],
            'COMPLAINT_NUMBER'   => (string) $data['complaintNumber'],
            'CREATOR_ID'         => (int) $data['creatorId'],
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => (int) $data['expertId'],
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => (int) $data['clientId'],
            'CLIENT_NAME'        => $client->getFullName(),
            'CONSUMPTION_ID'     => (int) $data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string) $data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string) $data['agreementNumber'],
            'DATE'               => (string) $data['date'],
            'DIFFERENCES'        => $differences,
        ];
    }

    private function validateTemplateData(array $templateData): void
    {
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }
    }

    private function sendEmployeeNotifications(string $emailFrom, array $templateData, int $resellerId): void
    {
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                $this->sendMessageToEmployee($emailFrom, $email, $templateData, $resellerId);
            }
        }
    }

    private function sendMessageToEmployee(string $emailFrom, string $email, array $templateData, int $resellerId): void
    {
        MessagesClient::sendMessage([
            0 => [
                'emailFrom' => $emailFrom,
                'emailTo'   => $email,
                'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
            ],
        ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
        $this->result['notificationClientByEmail'] = true;
    }

    private function sendClientNotifications(int $notificationType, array $data, string $emailFrom, Contractor $client, array $templateData, int $resellerId): void
    {
        // Send client email notification
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to']) && !empty($emailFrom) && !empty($client->email)) {
            $this->sendMessageToClient($emailFrom, $client->email, $templateData, $resellerId, $data['differences']['to']);
        }

        // Send client SMS notification
        $this->sendClientSmsNotification($notificationType, $data, $client, $templateData, $resellerId);
    }

    private function sendMessageToClient(string $emailFrom, string $clientEmail, array $templateData, int $resellerId, int $differenceTo): void
    {
        MessagesClient::sendMessage([
            0 => [
                'emailFrom' => $emailFrom,
                'emailTo'   => $clientEmail,
                'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
            ],
        ], $resellerId, $differenceTo);
        $this->result['notificationClientByEmail'] = true;
    }
    private function sendClientSmsNotification(int $notificationType, array $data, Contractor $client, array $templateData, int $resellerId): void
    {
        // Send client SMS notification
        if ($notificationType === self::TYPE_CHANGE && !empty($client->mobile)) {
            $error = '';
            $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);

            if ($res) {
                $this->result['notificationClientBySms']['isSent'] = true;
            }

            if (!empty($error)) {
                $this->result['notificationClientBySms']['message'] = $error;
            }
        }
    }
}
