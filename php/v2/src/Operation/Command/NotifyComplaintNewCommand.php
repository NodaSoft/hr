<?php

namespace NodaSoft\Operation\Command;

use NodaSoft\Messenger\Message;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\InitialData\NotifyComplaintNewInitialData;
use NodaSoft\Operation\Result\Result;
use NodaSoft\Operation\Result\NotifyComplaintNewResult;

class NotifyComplaintNewCommand implements Command
{
    /** @var NotifyComplaintNewResult */
    private $result;

    /** @var NotifyComplaintNewInitialData */
    private $initialData;

    /** @var Messenger */
    private $mail;

    /**
     * @param NotifyComplaintNewResult $result
     * @return void
     */
    public function setResult(Result $result): void
    {
        $this->result = $result;
    }

    public function setInitialData(InitialData $initialData): void
    {
        $this->initialData = $initialData;
    }

    public function setMail(Messenger $mail): void
    {
        $this->mail = $mail;
    }

    /**
     * @return NotifyComplaintNewResult
     */
    public function execute(): Result
    {
        $data = $this->initialData;
        $reseller = $data->getReseller();

        $message = new Message($data->getNotification(), $data->getMessageTemplate());

        foreach ($data->getEmployees() as $employee) {
            $result = $this->mail->send($message, $employee, $reseller);
            $this->result->addEmployeeEmailResult($result);
        }

        return $this->result;
    }
}
