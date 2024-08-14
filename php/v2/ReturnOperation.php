<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * Class ReturnOperation
 *
 * Обрабатывает операции возврата и отправляет уведомления.
 *
 * Переименовали на ReturnOperation чтобы соответсвовал требованиям PSR-4
 */
class ReturnOperation extends ReferencesOperation
{
    private const TYPE_NEW = 1;
    private const TYPE_CHANGE = 2;
    private const EVENT_TS_GOODS_RETURN = 'tsGoodsReturn'; // Константа для события
    private const SUCCESS_MESSAGE = 'success';

    /**
     * Выполняет операцию возврата и отправляет уведомления.
     *
     * @return array
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data = (array)$this->getRequest('data');
        $resellerId = isset($data['resellerId']) ? (int)$data['resellerId'] : 0; // Приведение к типу int
        $notificationType = (int)$data['notificationType'];
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        // Проверка наличия resellerId и его валидности
        if ($resellerId <= 0) {
            $result['notificationClientBySms']['message'] = 'Invalid or missing resellerId';
            return $result;
        }

        // Проверка наличия notificationType
        if ($notificationType <= 0) {
            throw new \Exception('Empty or invalid notificationType', 400);
        }

        // Получение реселлера
        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        // Получение клиента
        $client = Contractor::getById((int)$data['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Client not found or mismatch!', 400);
        }

        // Получение полного имени клиента
        $cFullName = $client->getFullName();
        if (empty($cFullName)) {
            $cFullName = $client->name; // Используем значение name, если getFullName возвращает пустую строку
        }

        // Получение данных о создателе и эксперте
        $cr = Employee::getById((int)$data['creatorId']);
        if ($cr === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $et = Employee::getById((int)$data['expertId']);
        if ($et === null) {
            throw new \Exception('Expert not found!', 400);
        }

        // Подготовка сообщения о различиях
        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO' => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }

        // Подготовка данных для шаблона
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

        // Проверка наличия всех данных для шаблона
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        // Получение email отправителя и списка email сотрудников
        $emailFrom = getResellerEmailFrom($resellerId);
        $emails = getEmailsByPermit($resellerId, self::EVENT_TS_GOODS_RETURN); // Использование константы для события

        // Отправка уведомлений сотрудникам
        $result['notificationEmployeeByEmail'] = $this->sendEmailNotifications($emails, $emailFrom, $templateData, $resellerId);

        // Отправка уведомлений клиентам
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                $this->sendEmailNotifications([$client->email], $emailFrom, $templateData, $resellerId, $client->id, (int)$data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                // Проверка наличия 'to' перед приведением к целому типу
                $to = !empty($data['differences']['to']) ? (int)$data['differences']['to'] : 0;
                // не понял смысла error и убрал, так как уведомление юзеру отправляем только при успешных
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $to, $templateData);
                // Смысла не было $res проверять и если true то присваивать true так как $res и так boolean
                $result['notificationClientBySms']['isSent'] = $res;
                $result['notificationClientBySms']['message'] = self::SUCCESS_MESSAGE;
            }
        }

        return $result;
    }

    //Во избежания дублирование кода добавлен новый метод
    /**
     * Отправляет уведомления по электронной почте.
     *
     * @param array $emails Массив email адресов получателей.
     * @param string $emailFrom Email отправителя.
     * @param array $templateData Данные для шаблона сообщения.
     * @param int $resellerId Идентификатор реселлера.
     * @param int $clientId Идентификатор клиента (опционально).
     * @param int $notificationType Тип уведомления (опционально).
     *
     * @return bool
     */
    private function sendEmailNotifications(array $emails, string $emailFrom, array $templateData, int $resellerId, int $clientId = 0, int $notificationType = 0): bool
    {
        if (empty($emailFrom) || empty($emails)) {
            return false;
        }

        foreach ($emails as $email) {
            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo' => $email,
                    'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, $clientId, $notificationType);
        }

        return true;
    }
}
