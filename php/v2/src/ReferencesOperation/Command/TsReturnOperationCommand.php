<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\Mail\Mail;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\MailFactory\TsReturnOperationComplaintMessageFactory;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;
use NodaSoft\ReferencesOperation\Result\TsReturnOperationResult;
use NW\WebService\References\Operations\Notification\NotificationEvents;
use NW\WebService\References\Operations\Notification\NotificationManager;

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
     * @var Mail
     */
    private $mail;

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

    public function setMail(Mail $mail): void
    {
        $this->mail = $mail;
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

        $messageFactory = new TsReturnOperationComplaintMessageFactory();

        foreach ($initialData->getEmployees() as $employee) {
            $message = $messageFactory->makeMessage($employee, $initialData);
            $result = $this->mail->send($message);
            $this->result->addEmployeeEmailResult($result);
        }

        $message = $messageFactory->makeMessage($client, $initialData);
        $result = $this->mail->send($message);
        $this->result->setClientEmailResult($result); //todo: handle logic

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE
            && ! is_null($initialData->getDifferencesTo())) {

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
