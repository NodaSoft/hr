<?php

namespace NW\WebService\References\Operations\Notification;

// название класса должно быть согласно названию файла для правильной работы autoloader
class ReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     */
    public function doOperation(): array // тип возвращаемого метода согласно интерфейсу
    {
        // потенциально небезопасное приведение неизвестного типа возвращаемого значения в массив
        $data = (array)$this->getRequest('data');

        $resellerId = filter_var($data['resellerId'] ?? null, FILTER_VALIDATE_INT);
        $notificationType = filter_var($data['notificationType'] ?? null, FILTER_VALIDATE_INT);

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

        // Проверка значения `reseller` не требуется, т.к. `Seller::getById()` всегда возвращает объект
        // Можно обработать потенциальное исключение конструктора (других исключений данный метод не вызывает)
        // По итогу переменная `$reseller` нигде не используется

        $clientId = filter_var($data['clientId'] ?? null, FILTER_VALIDATE_INT);
        if (empty($clientId)) {
            throw new \Exception('Empty clientId', 400);
        }

        $client = Contractor::getById($clientId);
        // проверка значения `client` на `null` не требуется, т.к. `Seller::getById()` всегда возвращает объект
        if ($client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Client not found!', 400);
        }

        // Здесь был лишний код ради кода, где метод `$client->getFullName()` вызывался 2 раза,
        // вместо проверки значения `$cFullName`

        // Судя по названию переменной `$cFullName` в ней ожидается ПОЛНОЕ имя клиента (name + id),
        // а не только имя (name), которое присваивалось переменной в случае если `getFullName()` возвращает пустую строку

        $cFullName = $client->getFullName(); // чисто теоретически, такое может быть согласно текущей реализации
        if (empty($cFullName)) {
            $cFullName = $client->name; // присваемое значение отличается от текущей реализации
            // Возможно, имелось в виду:
            // $cFullName = '$client->name' . ' ' . $client->id;
            // ?
        }

        // вот такие несколько строчек можно заменить одной если вынести в отдельный защищённый метод класса
        $creatorId = filter_var($data['creatorId'] ?? null, FILTER_VALIDATE_INT);
        if (empty($creatorId)) {
            throw new \Exception('Empty creatorId', 400);
        }

        $cr = Employee::getById($creatorId); // всегда возвращает объект

        // вот такие несколько строчек можно заменить одной если вынести в отдельный защищённый метод класса
        $expertId = filter_var($data['expertId'] ?? null, FILTER_VALIDATE_INT);
        if (empty($expertId)) {
            throw new \Exception('Empty expertId', 400);
        }

        $et = Employee::getById($expertId); // всегда возвращает объект

        $differences = $this->findDifferences($data, $notificationType, $resellerId);

        // Все приведения типов ниже являются небезопасными
        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'], // potential Warning: Object of class * could not be converted to int
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'], // potential Warning: Array to string conversion
            'CREATOR_ID'         => (int)$data['creatorId'], // potential Warning: Object of class * could not be converted to int
            'CREATOR_NAME'       => $cr->getFullName(),
            'EXPERT_ID'          => (int)$data['expertId'], // potential Warning: Object of class * could not be converted to int
            'EXPERT_NAME'        => $et->getFullName(),
            'CLIENT_ID'          => (int)$data['clientId'], // potential Warning: Object of class * could not be converted to int
            'CLIENT_NAME'        => $cFullName,
            'CONSUMPTION_ID'     => (int)$data['consumptionId'], // potential Warning: Object of class * could not be converted to int
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'], // potential Warning: Array to string conversion
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'], // potential Warning: Array to string conversion
            'DATE'               => (string)$data['date'], // potential Warning: Array to string conversion
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) { // такая проверка может работать неправильно в случае `false`, `0` или `[]`
                // вызов исключения !== "не отправляем уведомления"?
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId); // метод не принимает передаваемый аргумент
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (is_array($emails) && count($emails) > 0) { // `$emails && is_array($emails)`?
            foreach ($emails as $email) {
                // неизвестный класс `MessagesClient`
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $email,
                            // магическая функция не определена
                           'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                            // магическая функция не определена
                           'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $result['notificationEmployeeByEmail'] = true;

            }
        }

        // используется далее
        $differencesTo = filter_var($data['differences']['to'] ?? null, FILTER_VALIDATE_INT);

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($differencesTo)) {
            if (!empty($emailFrom) && !empty($client->email)) { // `$client` не имеет поля `email`
                // неизвестный класс `MessagesClient`
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $client->email,
                            // магическая функция не определена
                           'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                            // магическая функция не определена
                           'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) { // `$client` не имеет поля `mobile`
                // Undefined variable '$error'
                // Потенциально, ожидается возврат ошибки через передачу указателя на переменную `&$error`
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo, $templateData, $error);
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

    /**
     * @throws \Exception
     */
    protected function findDifferences(array $data, mixed $notificationType, mixed $resellerId): string
    {
        $differences = '';

        if ($notificationType === self::TYPE_NEW) {
            // магическая функция `__(string $action, ?array $attributes, int $resellerId)` не определена
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'] ?? null)) {
            $differencesFrom = filter_var($data['differences']['from'] ?? null, FILTER_VALIDATE_INT);
            if (empty($differencesFrom)) {
                throw new \Exception('Empty differences from', 400);
            }

            $differencesTo = filter_var($data['differences']['to'] ?? null, FILTER_VALIDATE_INT);
            if (empty($differencesTo)) {
                throw new \Exception('Empty differences to', 400);
            }

            // магическая функция `__(string $action, ?array $attributes, int $resellerId)` не определена
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName($differencesFrom),
                'TO'   => Status::getName($differencesTo),
            ], $resellerId);
        }

        return $differences;
    }
}
