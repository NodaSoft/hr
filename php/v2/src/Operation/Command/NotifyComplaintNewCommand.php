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
    /** @var NotifyComplaintNewInitialData */
    private $initialData;

    /** @var Messenger */
    private $mail;

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
        $result = new NotifyComplaintNewResult();
        $data = $this->initialData;
        $reseller = $data->getReseller();

        $message = new Message($data->getNotification(), $data->getMessageContentList());

        foreach ($data->getEmployees() as $employee) {
            $result->addEmployeeEmailResult(
                $this->mail->send($message, $employee, $reseller)
            );
        }

        return $result;
    }
}
