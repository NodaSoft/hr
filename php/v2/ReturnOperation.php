<?php

namespace NW\WebService\References\Operations\Notification;

use NodaSoft\Factory\OperationInitialData\TsReturnOperationInitialDataFactory;
use NodaSoft\Result\Operation\ReferencesOperation\ReferencesOperationResult;
use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /** @var TsReturnOperationResult */
    private $result;

    public function __construct()
    {
        $this->result = new TsReturnOperationResult();
    }

    /**
     * @throws \Exception
     * @return TsReturnOperationResult
     */
    public function doOperation(): ReferencesOperationResult
    {
        try {
            $params = $this->getRequest('data');
            $dataFactory = new TsReturnOperationInitialDataFactory();
            $initialData = $dataFactory->makeInitialData($params);
        } catch (\Exception $e) {
            //todo: handle an exception
            $this->result->setClientSmsErrorMessage($e->getMessage());
            return $this->result;
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
                $this->result->markEmployeeEmailSent();

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
                $this->result->markClientEmailSent();
            }

            if (!empty($client->mobile)) {
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
                if ($res) {
                    $this->result->markClientSmsSent();
                }
                if (!empty($error)) {
                    $this->result->setClientSmsErrorMessage($error);
                }
            }
        }

        return $this->result;
    }
}
