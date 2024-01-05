<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

class ReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws Exception
     */
    public function doOperation(): array
    {
        $data = new ReturnOperationRequestData($this->getRequest('data'));

        $resellerId = $data->resellerId;
        $notificationType = $data->notificationType;
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
            throw new Exception('Empty notificationType', 400);
        }

        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName($data->differences['from']),
                'TO' => Status::getName($data->differences['to']),
            ], $resellerId);
        } else {
            $differences = '';
        }

        if (!$differences) {
            throw new Exception("Differences in data is empty!", 500);
        }

        // Этот блок нигде не применяется. Возможно, он нужен только для проверки существования продавца
        try {
            $reseller = new Seller($resellerId);
        } catch (ExceptionAPI $e) {
            throw new Exception('Seller not found!', 400);
        } catch (Exception $e) {
            // Это не наша ошибка, прокидываем её как есть
            throw new Exception($e->getMessage(), $e->getCode());
        }

        try {
            $client = new Contractor($data->clientId);
            if ($client->type !== Contractor::TYPE_CUSTOMER || $client->seller->id !== $resellerId) {
                // Тут непонятные условия, но, не зная бизнес-логики и схемы данных, сложно сказать, что должно указывать на истинность клиента
                throw new Exception('Client not found!', 400);
            }
        } catch (ExceptionAPI $e) {
            throw new Exception('Client not found!', 400);
        } catch (Exception $e) {
            // Это не наша ошибка, прокидываем её как есть
            throw new Exception($e->getMessage(), $e->getCode());
        }

        $cFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $cFullName = $client->name;
        }

        try {
            $cr = new Employee($data->creatorId);
        } catch (ExceptionAPI $e) {
            throw new Exception('Creator not found!', 400);
        } catch (Exception $e) {
            // Это не наша ошибка, прокидываем её как есть
            throw new Exception($e->getMessage(), $e->getCode());
        }

        try {
            $et = new Employee($data->expertId);
        } catch (ExceptionAPI $e) {
            throw new Exception('Expert not found!', 400);
        } catch (Exception $e) {
            // Это не наша ошибка, прокидываем её как есть
            throw new Exception($e->getMessage(), $e->getCode());
        }

        $templateData = [
            'COMPLAINT_ID'       => $data->complaintId,
            'COMPLAINT_NUMBER'   => $data->complaintNumber,
            'CREATOR_ID'         => $data->creatorId,
            'CREATOR_NAME'       => $cr->getFullName(),
            'EXPERT_ID'          => $data->expertId,
            'EXPERT_NAME'        => $et->getFullName(),
            'CLIENT_ID'          => $data->clientId,
            'CLIENT_NAME'        => $cFullName,
            'CONSUMPTION_ID'     => $data->consumptionId,
            'CONSUMPTION_NUMBER' => $data->consumptionNumber,
            'AGREEMENT_NUMBER'   => $data->agreementNumber,
            'DATE'               => $data->date,
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = Config::getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = Config::getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            $messages = [];
            foreach ($emails as $email) {
                $messages[] = [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $email,
                    'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                ];
                $result['notificationEmployeeByEmail'] = true;
            }
            MessagesClient::sendMessage($messages, $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data->differences['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage(
                    [
                        [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $client->email,
                           'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                           'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                        ]
                    ],
                    $resellerId,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $client->id,
                    $data->differences['to']
                );
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $error = '';
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $data->differences['to'], $templateData, $error);
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
