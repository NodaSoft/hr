<?php

//declare cstrict type

namespace NW\WebService\References\Operations\Notification;

//final
// Ах да, не увидел еще одну реализацию ReferencesOperation, может лкчше сделать inteface чем абстрактный класс
class TsReturnOperation extends ReferencesOperation
{
    // Я бы консанты перенес в отдельный класс, и сделал как value object
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     */
    public function doOperation(): void //изменить тип
    {
        // Когда поулчаешь данные с http, 1ое, делаешь, сериализацию(мы же хотим в обьектам все держать), следом валидацию, как видим тут этого нет
        $data = (array)$this->getRequest('data');
        $resellerId = $data['resellerId'];
        $notificationType = (int)$data['notificationType'];

        // Dto нужно сделать
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        if (empty((int)$resellerId)) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $result;
        }

        //уже int есть выше
        if (empty((int)$notificationType)) {
            //Лучше делать ошибки более профильными. Если у тебя bundle/domain/модуль
            //final class OrderOutOfBoundsException extends \OutOfBoundsException
            //{
            //} И уже использвоать
            // Плюс тут просто exception нужно указать тип

            throw new \Exception('Empty notificationType', 400);
        }

        //Seller::getById - вообще странный метод, ты хочишь получить айди, а и тебе возвращается обьект класса. сомнительно и не окэй
        $reseller = Seller::getById((int)$resellerId);
        // Можно было задать нормальный констуктор в Seller или в базовом классе и тогда этого бы не понадобилось
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        $client = Contractor::getById((int)$data['clientId']);

        //Если был бы констукртор этого бы не понадобилось $client === null || $client->type !== Contractor::TYPE_CUSTOMER
        // $client->Seller->id !== $resellerId - по идеи вообще ничего не понимаю нахрена тут свойства Seller, сделать нормальный value object и просто вызвать гет и сравнить значения
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('сlient not found!', 400);
        }

        $cFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            // Кек типо если взять полное имя и оно пусто, то вхять имя, но полное имя строиться из имени
            $cFullName = $client->name;
        }

        // Можно было задать нормальный констуктор в Seller или в базовом классе и тогда этого бы не понадобилось
        $cr = Employee::getById((int)$data['creatorId']);
        if ($cr === null) {
            //Лучше делать ошибки более профильными.
            throw new \Exception('Creator not found!', 400);
        }

        // Можно было задать нормальный констуктор в Seller или в базовом классе и тогда этого бы не понадобилось
        $et = Employee::getById((int)$data['expertId']);
        if ($et === null) {
            //Лучше делать ошибки более профильными.
            throw new \Exception('Expert not found!', 400);
        }

        // Господи спаси и сохрани, что за такое __('NewPositionAdded') - Это типо так класс создается?
        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                    'FROM' => Status::getName((int)$data['differences']['from']),
                    'TO'   => Status::getName((int)$data['differences']['to']),
                ], $resellerId);
        }

        //DTO сделать и полям указать конкретный тип, и тогда не нужно городить массив и ниже не нужна првоерка на налл
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

        //DTO сделать и полям указать конкретный тип, и тогда не нужно городить массив и ниже не нужна првоерка на налл
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId); // Кеk это глобальная функция не принимает полей
        // Получаем email сотрудников из настроек
        // Зачем тут аргуементы класса?
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');

        // Сделать отдеьный класс или композит классов, который будет отправлять сособщения на почту или на мобилку
        //Про синтаксис тоже промолчу __('complaintEmployeeEmailSubject)
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
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
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $client->email,
                           'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                           'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
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
