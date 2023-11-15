<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\Messenger\Message;
use NodaSoft\Messenger\Messenger;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationNewInitialData;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;
use NodaSoft\ReferencesOperation\Result\ReturnOperationNewResult;

class ReturnOperationNewCommand implements ReferencesOperationCommand
{
    /** @var ReturnOperationNewResult */
    private $result;

    /** @var ReturnOperationNewInitialData */
    private $initialData;

    /** @var Messenger */
    private $mail;

    /**
     * @param ReturnOperationNewResult $result
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

    public function setMail(Messenger $mail): void
    {
        $this->mail = $mail;
    }

    /**
     * @return ReturnOperationNewResult
     */
    public function execute(): ReferencesOperationResult
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
