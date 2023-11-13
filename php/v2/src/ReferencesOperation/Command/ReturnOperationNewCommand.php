<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\Message\Messenger;
use NodaSoft\Message\Template\ComplaintNew;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\Message\Factory\MessageFactory;
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
        $initialData = $this->initialData;

        $complaintTemplate = new MessageFactory(new ComplaintNew());
        $reseller = $initialData->getReseller();

        foreach ($initialData->getEmployees() as $employee) {
            $message = $complaintTemplate->makeMessage($employee, $reseller, $initialData);
            $result = $this->mail->send($message);
            $this->result->addEmployeeEmailResult($result);
        }

        return $this->result;
    }
}
