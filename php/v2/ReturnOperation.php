<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends AbstractReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data = (array) $this->getRequest('data');
        $resellerId = (int) $data['resellerId'] ?? null;
        $notificationType = (int) $data['notificationType'] ?? null;
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
            throw new \Exception('Empty notificationType', 400);
        }

        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        $client = Contractor::getById((int) $data['clientId']);
        if ($client === null || $client->getType() !== Contractor::TYPE_CUSTOMER || $client->Seller->getId() !== $resellerId) {
            throw new \Exception('Client not found!', 400);
        }

        $cFullName = $client->getFullName() ?: $client->getName();

        $cr = Employee::getById((int)$data['creatorId']);
        if ($cr === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $et = Employee::getById((int)$data['expertId']);
        if ($et === null) {
            throw new \Exception('Expert not found!', 400);
        }

        $differences = match (true) {
            $notificationType === self::TYPE_NEW => __('NewPositionAdded', null, $resellerId),
            $notificationType === self::TYPE_CHANGE && !empty($data['differences']) => __(
                'PositionStatusHasChanged',
                [
                    'FROM' => StatusDictionary::getName((int)$data['differences']['from']),
                    'TO' => StatusDictionary::getName((int)$data['differences']['to']),
                ],
                $resellerId
            ),
            default => '',
        };

        $templateData = [
            'COMPLAINT_ID' => (int)$data['complaintId'],
            'COMPLAINT_NUMBER' => (string)$data['complaintNumber'],
            'CREATOR_ID' => (int)$data['creatorId'],
            'CREATOR_NAME' => $cr->getFullName(),
            'EXPERT_ID' => (int)$data['expertId'],
            'EXPERT_NAME' => $et->getFullName(),
            'CLIENT_ID' => (int)$data['clientId'],
            'CLIENT_NAME' => $cFullName,
            'CONSUMPTION_ID' => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER' => (string)$data['agreementNumber'],
            'DATE' => (string)$data['date'],
            'DIFFERENCES' => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ($key) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom();
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom)) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    [
                       'emailFrom' => $emailFrom,
                       'emailTo'   => $email,
                       'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                       'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEventsEnum::CHANGE_RETURN_STATUS->value);
                $result['notificationEmployeeByEmail'] = true;
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->getEmail())) {
                MessagesClient::sendMessage([
                    [
                       'emailFrom' => $emailFrom,
                       'emailTo'   => $client->getEmail(),
                       'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                       'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->getId(), NotificationEventsEnum::CHANGE_RETURN_STATUS->value, (int)$data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->getMobile())) {
                $res = NotificationManager::send(
                    $resellerId,
                    $client->getId(),
                    NotificationEventsEnum::CHANGE_RETURN_STATUS->value,
                    (int)$data['differences']['to'],
                    $templateData,
                    $error,
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
