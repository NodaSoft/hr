<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

class ReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        /**
         * Т.к. рефакторинг происходит в отрыве от остальной части
         * программы, менять тип входящих и исходящих данных нельзя.
         * В дальнейшем следует перейти от использования массивов
         * к использованию объектов DTO. Это обеспечит безопасную
         * передачу данных от метода к методу.
         *
         * @todo Входящие и выходящие данные сделать DTO
         */
        $data = (array)$this->getRequest('data');
        $resellerId = (int)$data['resellerId'];
        $clientId = (int)$data['clientId'];
        $creatorId = (int)$data['creatorId'];
        $expertId = (int)$data['expertId'];
        $differencesFrom = (int)$data['differences']['from'];
        $differencesTo = (int)$data['differences']['to'];
        $notificationType = (int)$data['notificationType'];
        $complaintId = (int)$data['complaintId'];
        $consumptionId = (int)$data['consumptionId'];
        $complaintNumber = (string)$data['complaintNumber'];
        $consumptionNumber = (string)$data['consumptionNumber'];
        $agreementNumber = (string)$data['agreementNumber'];
        $date = (string)$data['date'];

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

        $reseller = Seller::getById(($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        $client = Contractor::getById($clientId);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('сlient not found!', 400);
        }

        $cFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $cFullName = $client->name;
        }

        $cr = Employee::getById($creatorId);
        if ($cr === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $et = Employee::getById($expertId);
        if ($et === null) {
            throw new \Exception('Expert not found!', 400);
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && $differencesFrom !== $differencesTo) {
            $differences = __('PositionStatusHasChanged', [
                    'FROM' => Status::getName($differencesFrom),
                    'TO'   => Status::getName($differencesTo),
                ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID'       => $complaintId,
            'COMPLAINT_NUMBER'   => $complaintNumber,
            'CREATOR_ID'         => $creatorId,
            'CREATOR_NAME'       => $cr->getFullName(),
            'EXPERT_ID'          => $expertId,
            'EXPERT_NAME'        => $et->getFullName(),
            'CLIENT_ID'          => $clientId,
            'CLIENT_NAME'        => $cFullName,
            'CONSUMPTION_ID'     => $consumptionId,
            'CONSUMPTION_NUMBER' => $consumptionNumber,
            'AGREEMENT_NUMBER'   => $agreementNumber,
            'DATE'               => $date,
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
        if ($notificationType === self::TYPE_CHANGE && !empty($differencesTo)) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $client->email,
                           'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                           'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
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
}
