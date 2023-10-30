<?php

namespace NW\WebService\References\Operations\Notification;

use NodaSoft\Factory\OperationInitialData\TsReturnOperationInitialDataFactory;
use NodaSoft\OperationParams\TsReturnOperationParams;
use NodaSoft\Request\HttpRequest;
use NodaSoft\Result\Operation\ReferencesOperation\ReferencesOperationResult;
use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     * @return TsReturnOperationResult
     */
    public function doOperation(): ReferencesOperationResult
    {
        $result = new TsReturnOperationResult();

        $params = new TsReturnOperationParams();
        $params->setRequest(new HttpRequest());
        if ($params->isValid()) {
            $result->setClientSmsErrorMessage('Required parameter is missing.');
            return $result;
        }

        try {
            $dataFactory = new TsReturnOperationInitialDataFactory();
            $initialData = $dataFactory->makeInitialData($params);
        } catch (\Exception $e) {
            //todo: handle an exception
            $result->setClientSmsErrorMessage($e->getMessage());
            return $result;
        }

        $templateData = $initialData->getMessageTemplate()->toArray();
        $resellerId = $initialData->getReseller();
        $notificationType = $initialData->getNotificationType();

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
                $result->markEmployeeEmailSent();

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
                $result->markClientEmailSent();
            }

            if (!empty($client->mobile)) {
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
                if ($res) {
                    $result->markClientSmsSent();
                }
                if (!empty($error)) {
                    $result->setClientSmsErrorMessage($error);
                }
            }
        }

        return $result;
    }
}
