<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws Exception
     */
    public function doOperation(): array
    {
        $data = (array)$this->getRequest('data');


        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];
        $resellerId = (int)$data['resellerId'];
        if (!$resellerId) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $result;
        }
        $notificationType = (int)$data['notificationType'];
        if (!$notificationType) {
            throw new Exception('Empty notificationType', 400);
        }

        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new Exception('Seller not found!', 400);
        }

        $client = Contractor::getById((int)$data['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
         throw new Exception('Client not found or not valid!', 400);
        }

        $clientFullName = $client->getFullName() ?: $client->name;

        $creator = Employee::getById((int)$data['creatorId']);
        if ($creator === null) {
            throw new Exception('Creator not found!', 400);
        }

        $expert = Employee::getById((int)$data['expertId']);
        if ($expert === null) {
            throw new Exception('Expert not found!', 400);
        }

        $differences = $this->getDifferencesMessage($notificationType, $data, $resellerId);

        $templateData = $this->prepareTemplateData($data, $clientFullName, $creator, $expert, $differences);

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $value) {
            if (empty($value)) {
                throw new Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId);
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');

        if ($emailFrom && count($emails) > 0) {
            $this->sendEmployeeNotifications($emails, $emailFrom, $templateData, $resellerId, $result);
        }

        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            $this->sendClientNotifications($emailFrom, $client, $resellerId, $data, $templateData, $result);
        }

        return $result;
    }

    private function getDifferencesMessage(int $notificationType, array $data, int $resellerId): string
    {
        if ($notificationType === self::TYPE_NEW) {
            return __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            return __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO'   => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }
        return '';
    }

    private function prepareTemplateData(array $data, string $clientFullName, $creator, $expert, string $differences): array
    {
        return [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => (int)$data['creatorId'],
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => (int)$data['expertId'],
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => (int)$data['clientId'],
            'CLIENT_NAME'        => $clientFullName,
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $differences,
        ];
    }

    private function sendEmployeeNotifications(array $emails, string $emailFrom, array $templateData, int $resellerId, array &$result): void
    {
        foreach ($emails as $email) {
            MessagesClient::sendMessage([
                0 => [
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $email,
                    'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
            $result['notificationEmployeeByEmail'] = true;
        }
    }

    private function sendClientNotifications(string $emailFrom, $client, int $resellerId, array $data, array $templateData, array &$result): void
    {
        if ($emailFrom && !empty($client->email)) {
            MessagesClient::sendMessage([
                0 => [
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $client->email,
                    'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                    'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);
            $result['notificationClientByEmail'] = true;
        }

        if (!empty($client->mobile)) {
            $error = '';
            $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
            if ($res) {
                $result['notificationClientBySms']['isSent'] = true;
            }
            if ($error) {
                $result['notificationClientBySms']['message'] = $error;
            }
        }
    }
}
