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
        $data = $this->getRequest('data') ? $this->getRequest('data') : null;

        if (is_array($data) && !empty($data['differences']) && is_array($data['differences'])) {
            $resellerId = is_numeric($data['resellerId']) ? (int)$data['resellerId'] : null;
            $clientId = is_numeric($data['clientId']) ? (int)$data['clientId'] : null;
            $notificationType = is_numeric($data['notificationType']) ? (int)$data['notificationType'] : null;
            $creatorId = is_numeric($data['creatorId']) ? (int)$data['creatorId'] : null;
            $expertId = is_numeric($data['expertId']) ? (int)$data['expertId'] : null;
            $complaintId = is_numeric($data['complaintId']) ? (int)$data['complaintId'] : null;
            $consumptionId = is_numeric($data['consumptionId']) ? (int)$data['consumptionId'] : null;
        } else {
            throw new Exception('Invalid request', 400);
        }

        if (empty($resellerId)) {
            throw new Exception('Empty resellerId', 400);
        }

        if (empty($clientId)) {
            throw new Exception('Empty clientId', 400);
        }

        if (empty($notificationType)) {
            throw new Exception('Empty notificationType', 400);
        }

        if (empty($creatorId)) {
            throw new Exception('Empty creatorId', 400);
        }

        if (empty($expertId)) {
            throw new Exception('Empty expertId', 400);
        }

        $reseller = Seller::getById($resellerId);
        if (empty($reseller)) {
            throw new Exception('Seller not found!', 400);
        }

        $client = Contractor::getById($clientId);
        if (empty($client) || $client->type !== Contractor::TYPE_CUSTOMER || $client->seller->id !== $resellerId) {
            throw new Exception('Client not found!', 400);
        }

        $clientFullName = trim($client->getFullName()); // может быть пустым только если $client->name пустой
        /*if (empty($client->getFullName())) {
            $clientFullName = $client->name;
        }*/

        $creator = Employee::getById($creatorId);
        if (empty($creator)) {
            throw new Exception('Creator not found!', 400);
        }

        $expert = Employee::getById($expertId);
        if (empty($expert)) {
            throw new Exception('Expert not found!', 400);
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])
            && !empty($data['differences']['from'])) {
            $differences = __('PositionStatusHasChanged', [
                    'FROM' => Status::getName((int)$data['differences']['from']),
                    'TO'   => Status::getName((int)$data['differences']['to']),
                ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID'       => $complaintId,
            'COMPLAINT_NUMBER'   => $data['complaintNumber'], // нет смысла привести к типу string
            'CREATOR_ID'         => $creatorId,
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $expertId,
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $clientId,
            'CLIENT_NAME'        => $clientFullName,
            'CONSUMPTION_ID'     => $consumptionId,
            'CONSUMPTION_NUMBER' => $data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => $data['agreementNumber'],
            'DATE'               => $data['date'],
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        $emailFrom = getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && !empty($emails)) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage(
                    [
                        (int)MessageTypes::EMAIL => [ // MessageTypes::EMAIL
                               'emailFrom' => $emailFrom,
                               'emailTo'   => $email,
                               'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                               'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                        ],
                    ],
                    $resellerId,
                    NotificationEvents::CHANGE_RETURN_STATUS
                );
                $result['notificationEmployeeByEmail'] = true;

            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage(
                    [
                        (int)MessageTypes::EMAIL => [ // MessageTypes::EMAIL
                               'emailFrom' => $emailFrom,
                               'emailTo'   => $client->email,
                               'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                               'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                        ],
                    ],
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    (int)$data['differences']['to']
                );
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $error = null; // возможно, меняется по ссылке

                $res = NotificationManager::send(
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    (int)$data['differences']['to'],
                    $templateData,
                    $error
                );
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
