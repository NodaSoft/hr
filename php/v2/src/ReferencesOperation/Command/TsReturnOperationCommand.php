<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;
use NodaSoft\ReferencesOperation\Result\TsReturnOperationResult;
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

    /**
     * @var InitialData
     */
    private $initialData;

    /**
     * @param TsReturnOperationResult $result
     * @return void
     */
    public function setResult(ReferencesOperationResult $result): void
    {
        $this->result = $result;
    }

    public function setInitialData(InitialData $initialData): void
    {
        $this->initialData = $initialData;
    }

    /**
     * @return TsReturnOperationResult
     */
    public function execute(): ReferencesOperationResult
    {
        $initialData = $this->initialData;
        $client = $initialData->getClient();
        $templateData = $initialData->getMessageTemplate()->toArray();
        $resellerId = $initialData->getReseller();
        $notificationType = $initialData->getNotificationType();

        $emailFrom = getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (! empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage(
                    [
                        0 => [ // MessageTypes::EMAIL
                            'emailFrom' => $emailFrom,
                            'emailTo'   => $email,
                            'subject'   => __(
                                'complaintEmployeeEmailSubject',
                                $templateData,
                                $resellerId
                            ),
                            'message'   => __(
                                'complaintEmployeeEmailBody',
                                $templateData,
                                $resellerId
                            ),
                        ],
                    ],
                    $resellerId,
                    NotificationEvents::CHANGE_RETURN_STATUS
                );
            }
            $this->result->markEmployeeEmailSent();
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE
            && ! is_null($initialData->getDifferencesTo())) {
            if (! empty($emailFrom) && ! empty($client->getEmail())) {
                MessagesClient::sendMessage(
                    [
                        0 => [ // MessageTypes::EMAIL
                            'emailFrom' => $emailFrom,
                            'emailTo'   => $client->getEmail(),
                            'subject'   => __(
                                'complaintClientEmailSubject',
                                $templateData,
                                $resellerId
                            ),
                            'message'   => __(
                                'complaintClientEmailBody',
                                $templateData,
                                $resellerId
                            ),
                        ],
                    ],
                    $resellerId,
                    $client->getId(),
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $initialData->getDifferencesTo()
                );
                $this->result->markClientEmailSent();
            }

            if (! empty($client->getCellphoneNumber())) {
                $result = NotificationManager::send(
                    $resellerId,
                    $client->getId(),
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $initialData->getDifferencesTo(),
                    $templateData
                );
                if ($result->hasError()) {
                    $this->result->setClientSmsErrorMessage($result->getErrorMessage());
                    return $this->result;
                }
                $this->result->markClientSmsSent();
            }
        }

        return $this->result;
    }
}
