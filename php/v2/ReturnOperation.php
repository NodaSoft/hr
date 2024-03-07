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
        $data = $this->getRequest('data');

        $notificationType = (int)$this->getValueFromArray($data,'notificationType');
        $clientId =(int)$this->getValueFromArray($data,'clientId');
        $creatorId = (int)$this->getValueFromArray($data,'creatorId');
        $expertId = (int)$this->getValueFromArray($data,'expertId');

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
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

        $client = Client::getById($clientId);
        if (!$client->isCustomer() || !$client->Seller->isIdEqualTo($reseller->id)) {
            throw new \Exception('Client not found!', 400);
        }

        $cFullName = $client->getFullName();

        $cr = Employee::getById($creatorId);

        $et = Expert::getById($expertId);

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($this->getValueFromArray($data,'differences'))) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$this->getValueFromArray($data,'differences','from')),
                'TO'   => Status::getName((int)$this->getValueFromArray($data,'differences','to')),
            ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID'       => (int)$this->getValueFromArray($data,'complaintId'),
            'COMPLAINT_NUMBER'   => (string)$this->getValueFromArray($data,'complaintNumber'),
            'CREATOR_ID'         => (int)$this->getValueFromArray($data,'creatorId'),
            'CREATOR_NAME'       => $cr->getFullName(),
            'EXPERT_ID'          => (int)$this->getValueFromArray($data,'expertId'),
            'EXPERT_NAME'        => $et->getFullName(),
            'CLIENT_ID'          => (int)$this->getValueFromArray($data,'clientId'),
            'CLIENT_NAME'        => $cFullName,
            'CONSUMPTION_ID'     => (int)$this->getValueFromArray($data,'consumptionId'),
            'CONSUMPTION_NUMBER' => (string)$this->getValueFromArray($data,'consumptionNumber'),
            'AGREEMENT_NUMBER'   => (string)$this->getValueFromArray($data,'agreementNumber'),
            'DATE'               => (string)$this->getValueFromArray($data,'date'),
            'DIFFERENCES'        => $differences,
        ];

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
                try {
                    MessagesClient::sendMessage([
                        0 => [ // MessageTypes::EMAIL
                            'emailFrom' => $emailFrom,
                            'emailTo'   => $email,
                            'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                            'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                        ],
                    ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                } catch (\Exception $exception) {

                }
                $result['notificationEmployeeByEmail'] = true;

            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($this->getValueFromArray($data,'differences','to'))) {
            if (!empty($emailFrom) && !empty($client->email)) {
                try {
                    MessagesClient::sendMessage([
                        0 => [ // MessageTypes::EMAIL
                            'emailFrom' => $emailFrom,
                            'emailTo' => $client->email,
                            'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                            'message' => __('complaintClientEmailBody', $templateData, $resellerId),
                        ],
                    ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$this->getValueFromArray($data, 'differences', 'to'));
                } catch (\Exception $exception) {

                }
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $error=null;
                $res=null;
                try {
                    $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$this->getValueFromArray($data, 'differences', 'to'), $templateData, $error);
                } catch (\Exception $exception) {
                    $error='Unknown error';
                }
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

    public function getValueFromArray(array $data,...$keys)
    {
        foreach ($keys as $key) {
            $data=$data[$key]??null;
        }
        return $data;
    }
}
