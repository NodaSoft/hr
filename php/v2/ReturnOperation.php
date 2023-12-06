<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Clients\Seller;
use NW\WebService\References\Operations\Notification\Clients\Contractor;
use NW\WebService\References\Operations\Notification\Clients\Employee;
use NW\WebService\References\Operations\Notification\Enums\ClientType;
use NW\WebService\References\Operations\Notification\Enums\NotificationEvents;
use NW\WebService\References\Operations\Notification\Enums\NotificationType;
use NW\WebService\References\Operations\Notification\Enums\Status;

class ReturnOperation extends ReferencesOperation
{
    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data = (array) $this->getRequest('data');

        $resellerId = isset($data['resellerId'])
            ? (int) $data['resellerId']
            : null;

        $notificationType = isset($data['notificationType'])
            ? (int) $data['notificationType']
            : null;

        if (empty($resellerId)) {
            throw new \Exception('Empty resellerId', 400);
        }

        if (empty($notificationType)) {
            throw new \Exception('Empty notificationType', 400);
        }

        $reseller = Seller::getById($resellerId);
        if (empty($reseller)) {
            throw new \Exception('Seller not found!', 404);
        }

        $client = Contractor::getById((int)$data['clientId']);
        if (
            empty($client)
            || $client->type !== ClientType::CUSTOMER
//            || $client->Seller->id !== $resellerId // TODO: Неочень понятно что это за проверка
        ) {
            throw new \Exception('Client not found!', 404);
        }

        $creator = Employee::getById((int)$data['creatorId']);
        if (empty($creator)) {
            throw new \Exception('Creator not found!', 404);
        }

        $expert = Employee::getById((int)$data['expertId']);
        if (empty($expert)) {
            throw new \Exception('Expert not found!', 404);
        }

        $differences = '';
        if ($notificationType === NotificationType::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === NotificationType::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO' => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => $creator->id,
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $expert->id,
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $client->id,
            'CLIENT_NAME'        => $client->getFullName(),
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $value) {
            if (empty($value)) {
                throw new \Exception(sprintf('Template Data (%s) is empty!', $key), 400);
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

        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $email,
                        'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);

                $result['notificationEmployeeByEmail'] = true;
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === NotificationType::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    [ // MessageTypes::EMAIL
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
