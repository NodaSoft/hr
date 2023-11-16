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
    /** @var NotifyComplaintStatusChangedInitialData */
    private $initialData;

    /** @var Messenger */
    private $mail;

    /** @var Messenger */
    private $sms;

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
        $result = new ReturnOperationStatusChangedResult();
        $data = $this->initialData;
        $reseller = $data->getReseller();
        $client = $data->getClient();

        $message = new Message($data->getNotification(), $data->getMessageContentList());

        foreach ($data->getEmployees() as $employee) {
            $result->addEmployeeEmailResult($this->mail->send($message, $employee, $reseller));
        }

        $result->setClientEmailResult($this->mail->send($message, $client, $reseller));
        $result->setClientSmsResult($this->sms->send($message, $client, $reseller));

        return $result;
    }
}
