<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    /**
     * Выполняет операцию по возврату товаров и отправке уведомлений о изменениях статусов возвратов.
     *
     * @throws Exception
     */
    public function doOperation(): array
    {
        $data = $this->getRequest('data');
        $resellerId = $data['resellerId'] ?? null;
        $notificationType = $data['notificationType'] ?? null;
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        if (empty($resellerId)) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $result;
        }

        if (empty($notificationType)) {
            throw new Exception('Empty notificationType', 400);
        }

        $reseller = Seller::getById($resellerId);
        if (!$reseller) {
            throw new Exception('Seller not found!', 400);
        }

        $clientId = $data['clientId'] ?? null;
        $client = Contractor::getById($clientId);
        if (!$client || $client->type !== Contractor::TYPE_CUSTOMER || $client->id !== $resellerId) {
            throw new Exception('Client not found!', 400);
        }

        $creatorId = $data['creatorId'] ?? null;
        $cr = Employee::getById($creatorId);
        if (!$cr) {
            throw new Exception('Creator not found!', 400);
        }

        $expertId = $data['expertId'] ?? null;
        $et = Employee::getById($expertId);
        if (!$et) {
            throw new Exception('Expert not found!', 400);
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName($data['differences']['from'] ?? 0),
                'TO' => Status::getName($data['differences']['to'] ?? 0),
            ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID' => $data['complaintId'] ?? null,
            'COMPLAINT_NUMBER' => $data['complaintNumber'] ?? null,
            'CREATOR_ID' => $creatorId,
            'CREATOR_NAME' => $cr->getFullName(),
            'EXPERT_ID' => $expertId,
            'EXPERT_NAME' => $et->getFullName(),
            'CLIENT_ID' => $clientId,
            'CLIENT_NAME' => $client->getFullName() ?: $client->name,
            'CONSUMPTION_ID' => $data['consumptionId'] ?? null,
            'CONSUMPTION_NUMBER' => $data['consumptionNumber'] ?? null,
            'AGREEMENT_NUMBER' => $data['agreementNumber'] ?? null,
            'DATE' => $data['date'] ?? null,
            'DIFFERENCES' => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $tempData) {
            if (empty($tempData)) {
                throw new Exception('Template Data is empty!', 500);
            }
        }

        $emailFrom = getResellerEmailFrom();
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo' => $email,
                        'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $result['notificationEmployeeByEmail'] = true;
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo' => $client->email,
                        'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $clientId, NotificationEvents::CHANGE_RETURN_STATUS, $data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $res = NotificationManager::send($resellerId, $clientId, NotificationEvents::CHANGE_RETURN_STATUS, $data['differences']['to'], $templateData, $error);
                if ($res) {
                    $result['notificationClientBySms']['isSent'] = true;
                }
                if (!empty($error)) {
                    $result['notificationClientBySms']['message'] = $error;
                }
            }
        }

        return $result;
    }
}
