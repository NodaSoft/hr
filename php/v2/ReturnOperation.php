<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     */
    public function doOperation(): array // в абстрактном классе от которого унаследуем метдо должен вдавать array
    {
        $data = (array)$this->getRequest('data');

        // не факт что ключ "data" присутствует в теле $_REQUEST
        if (empty($data)) {
            throw new \Exception("Empty request data", 400);
        }

        // перемещаю инициализацию переменных $resellerId, $notificationType
        // непосредственно где их используют, так читаемость больше

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        $resellerId = $data['resellerId'];
        // на нужен (int) каст сразу внутри empty() = warning
        if (empty($resellerId) || !is_numeric($resellerId)) {
            $result['notificationClientBySms']['message'] = 'Empty or not valid numeric resellerId';
            return $result;
        }

        $reseller = Seller::getById((int)$resellerId);
        // (*) - смотреть в readme.md
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        // (int) каст на непроверенное значение массива ... смело! ... возможен Warning!
        // $notificationType = (int)$data['notificationType'];
        $notificationType = $data['notificationType'];  // лучше так
        if (empty($notificationType) || !is_numeric($notificationType)) {
            throw new \Exception('Empty or not valid numeric notificationType', 400);
        }

        $clientId = $data['clientId'];
        if (empty($clientId) || !is_numeric($clientId)) {
            throw new \Exception('Empty or not valid numeric clientId', 400);
        }
        $client = Contractor::getById((int)$clientId);
        // (*) - смотреть в readme.md
        // ... || $client->Seller->id !== $resellerId) {
        // Трюк с декларацией "@property Seller $Seller" в комментах а-ля Annotation, которая обманывает проверку
        // синтаксиса IDE (таких как phpStorm) не прокатит на run-time!
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER) {
            throw new \Exception('Client not found!', 400);
        }

        $cFullName = $client->getFullName();
        // (***) - смотреть в $templeteData
        // if (empty($client->getFullName())) {     как минимун $client->id всегда будет,
        //     $cFullName = $client->name;          поэтому $cFullName != '' всегда
        // }

        $creatorId = $data['creatorId'];
        if (empty($creatorId) || !is_numeric($creatorId)) {
            throw new \Exception('Empty or not valid numeric creatorId', 400);
        }
        $cr = Employee::getById((int)$creatorId);
        // (*) - смотреть в readme.md
        if ($cr === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $expertId = $data['expertId'];
        if (empty($expertId) || !is_numeric($expertId)) {
            throw new \Exception('Empty or not valid numeric expertId', 400);
        }
        $et = Employee::getById((int)$expertId);
        // (*) - смотреть в readme.md
        if ($et === null) {
            throw new \Exception('Expert not found!', 400);
        }



        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            // (**) смотреть в readme.md
            $differences = __('NewPositionAdded', null, $resellerId);
        }
        elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            // (**) смотреть в readme.md
            $differences = __('PositionStatusHasChanged', [
                    'FROM' => Status::getName((int)$data['differences']['from']),
                    'TO'   => Status::getName((int)$data['differences']['to']),
                ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => $creatorId,
            'CREATOR_NAME'       => $cr->getFullName(),     // (***) забавно что $cFullName проверяли а тут прям в лоб
            'EXPERT_ID'          => $expertId,
            'EXPERT_NAME'        => $et->getFullName(),     // (***) забавно что $cFullName проверяли а тут прям в лоб
            'CLIENT_ID'          => $clientId,
            'CLIENT_NAME'        => $cFullName,             // (***)
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
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

                // здесь я буду подразумевать что такой MessagesClient где-то в приложении существует!
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $email,
                            // (**) смотреть в readme.md
                           'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                           'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $result['notificationEmployeeByEmail'] = true;
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {

                // здесь я буду подразумевать что такой MessagesClient где-то в приложении существует!
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $client->email,
                            // (**) смотреть в readme.md
                           'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                           'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            // "$client->mobile" этой property в декларации класса нигде не было !!!.
            // либо там недочёт, либо это код из другого приложения (или другой версии приложения где есть новая фича).
            if (!empty($client->mobile)) {
                // здесь я буду подразумевать что такой NotificationManager где-то существует!
                // откуда переменная $error ?
                // $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
                // Можно конечно заранее её обозначить и передать ByReference, вдруг там что-то случится и тогда можно её наполнить смычслом!
                // $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, &$error);
                // Плохая имплементация NotificationManager.
                //
                // Я бы отдавал в $res = ['status' => OK | ERROR, 'error' => string ]
                // а в IF уже проверял статус $res!
                $error = "";
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, &$error);
                if ($res['status'] == 'ERROR') {
                    $result['notificationClientBySms']['message'] = $res['error'];
                } else {
                    $result['notificationClientBySms']['isSent'] = true;
                }
            }
        }

        return $result;
    }
}

