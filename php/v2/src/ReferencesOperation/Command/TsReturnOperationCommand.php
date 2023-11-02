<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\Factory\OperationInitialData\TsReturnOperationInitialDataFactory;
use NodaSoft\OperationParams\ReferencesOperationParams;
use NodaSoft\OperationParams\TsReturnOperationParams;
use NodaSoft\Result\Operation\ReferencesOperation\ReferencesOperationResult;
use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;
use NW\WebService\References\Operations\Notification\MessagesClient;
use NW\WebService\References\Operations\Notification\NotificationEvents;
use NW\WebService\References\Operations\Notification\NotificationManager;
use function NW\WebService\References\Operations\Notification\__;
use function NW\WebService\References\Operations\Notification\getEmailsByPermit;
use function NW\WebService\References\Operations\Notification\getResellerEmailFrom;

class TsReturnOperationCommand implements ReferencesOperationCommand
{
    public const TYPE_NEW = 1;

    public const TYPE_CHANGE = 2;

    /** @var TsReturnOperationResult */
    private $result;

    /** @var TsReturnOperationParams */
    private $params;

    /**
     * @param TsReturnOperationResult $result
     * @return void
     */
    public function setResult(ReferencesOperationResult $result): void
    {
        $this->result = $result;
    }

    /**
     * @param TsReturnOperationParams $params
     * @return void
     */
    public function setParams(ReferencesOperationParams $params): void
    {
        $this->params = $params;
    }

    /**
     * @return TsReturnOperationResult
     */
    public function execute(): ReferencesOperationResult
    {
        if ($this->params->isValid()) {
            $this->result->setClientSmsErrorMessage('Required parameter is missing.');
            return $this->result;
        }

        try {
            $dataFactory = new TsReturnOperationInitialDataFactory();
            $initialData = $dataFactory->makeInitialData($this->params);
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
        if (! empty($emailFrom) && count($emails) > 0) {
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
        if ($notificationType === self::TYPE_CHANGE && ! empty($initialData->dateTo())) {
            if (! empty($emailFrom) && ! empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $client->email,
                        'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $initialData->dateTo());
                $this->result->markClientEmailSent();
            }

            if (! empty($client->mobile)) {
                $res = NotificationManager::send(
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $initialData->dateTo(),
                    $templateData,
                    $error
                );
                if ($res) {
                    $this->result->markClientSmsSent();
                }
                if (! empty($error)) {
                    $this->result->setClientSmsErrorMessage($error);
                }
            }
        }

        return $this->result;
    }
}
