<?php

namespace NodaSoft\Operation\Command;

use NodaSoft\Messenger\Message;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\InitialData\NotifyComplaintStatusChangedInitialData;
use NodaSoft\Operation\Result\Result;
use NodaSoft\Operation\Result\ReturnOperationStatusChangedResult;

class NotifyComplaintStatusChangedCommand implements Command
{
    /** @var ReturnOperationStatusChangedResult */
    private $result;

    /** @var NotifyComplaintStatusChangedInitialData */
    private $initialData;

    /** @var Messenger */
    private $mail;

    /** @var Messenger */
    private $sms;

    /**
     * @param ReturnOperationStatusChangedResult $result
     * @return void
     */
    public function setResult(Result $result): void
    {
        $this->result = $result;
    }

    /**
     * @param NotifyComplaintStatusChangedInitialData $initialData
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
    public function execute(): Result
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
