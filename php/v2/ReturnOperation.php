<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
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
        $notificationType = null;
        if(isset($data['notificationType'])) {
            $notificationType = (int)$data['notificationType'];
        }

        if (empty((int)$notificationType)) {
            throw new \Exception('Empty notificationType', 400);
        }

        $resellerId = null;
        if(isset($data['resellerId'])) {
            $resellerId = $data['resellerId'];
            $reseller = Seller::getById((int)$resellerId);
        }
        if (empty((int)$resellerId)) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $result;
        }

        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        if(isset($data['clientId'])) {
            $client = Contractor::getById((int)$data['clientId']);
        }

        if ($client === null || $client->getType() !== Contractor::TYPE_CUSTOMER ||
        $client->Seller->getId() !== $resellerId) {
            throw new \Exception('client not found!', 400);
        }

        $cFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $cFullName = $client->getName();
        }

        $cr = null;
        if(isset($data['creatorId'])) {
            $cr = Employee::getById((int)$data['creatorId']);
        }
        if ($cr === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $et = null;
        if(isset($data['expertId'])) {
            $et = Employee::getById((int)$data['expertId']);
        }
        if ($et === null) {
            throw new \Exception('Expert not found!', 400);
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                    'FROM' => Status::getName((int)$data['differences']['from']),
                    'TO'   => Status::getName((int)$data['differences']['to']),
                ], $resellerId);
        }

        try {
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
        } catch (\Throwable $th) {
            throw new \Exception('Data not set', 400);
        }

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $email,
                           'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                           'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS->value);
                $result['notificationEmployeeByEmail'] = true;
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if(isset($data['differences']['to'])) {
            if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
                if (!empty($emailFrom) && !empty($client->getEmail())) {
                    MessagesClient::sendMessage([
                        0 => [ // MessageTypes::EMAIL
                               'emailFrom' => $emailFrom,
                               'emailTo'   => $client->getEmail(),
                               'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                               'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                        ],
                    ], $resellerId, $client->getId(), NotificationEvents::CHANGE_RETURN_STATUS->value, (int)$data['differences']['to']);
                    $result['notificationClientByEmail'] = true;
                }

                if (!empty($client->getMobile())) {
                    $error = 'Mobile didn`t set';
                    $res = NotificationManager::send($resellerId, $client->getId(), NotificationEvents::CHANGE_RETURN_STATUS->value, (int)$data['differences']['to'], $templateData, $error);
                    if ($res) {
                        $result['notificationClientBySms']['isSent'] = true;
                    }
                    if (!empty($error)) {
                        $result['notificationClientBySms']['message'] = $error;
                    }
                }
            }
        }

        return $result;
    }
}
