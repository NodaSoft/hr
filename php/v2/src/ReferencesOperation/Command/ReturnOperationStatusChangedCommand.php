<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\Messenger\Message;
use NodaSoft\Messenger\Messenger;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationStatusChangedInitialData;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;
use NodaSoft\ReferencesOperation\Result\ReturnOperationStatusChangedResult;

class ReturnOperationStatusChangedCommand implements ReferencesOperationCommand
{
    /** @var ReturnOperationStatusChangedResult */
    private $result;

    /** @var ReturnOperationStatusChangedInitialData */
    private $initialData;

    /** @var Messenger */
    private $mail;

    /** @var Messenger */
    private $sms;

    /**
     * @param ReturnOperationStatusChangedResult $result
     * @return void
     */
    public function setResult(ReferencesOperationResult $result): void
    {
        $this->result = $result;
    }

    /**
     * @param ReturnOperationStatusChangedInitialData $initialData
     * @return void
     */
    public function setInitialData(InitialData $initialData): void
    {
        $this->initialData = $initialData;
    }

    public function setMail(Messenger $mail): void
    {
        $this->mail = $mail;
    }

    public function setSms(Messenger $sms): void
    {
        $this->sms = $sms;
    }

    /**
     * @return ReturnOperationStatusChangedResult
     */
    public function execute(): ReferencesOperationResult
    {
        $data = $this->initialData;
        $reseller = $data->getReseller();
        $client = $data->getClient();

        $message = new Message($data->getNotification(), $data->getMessageTemplate());

        foreach ($data->getEmployees() as $employee) {
            $this->result->addEmployeeEmailResult($this->mail->send($message, $employee, $reseller));
        }

        $this->result->setClientEmailResult($this->mail->send($message, $client, $reseller));
        $this->result->setClientSmsResult($this->sms->send($message, $client, $reseller));

        return $this->result;
    }
}
