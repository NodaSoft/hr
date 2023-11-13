<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\Message\Messenger;
use NodaSoft\Message\Template\ComplaintStatusChanged;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\Message\Factory\MessageFactory;
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
        $initialData = $this->initialData;
        $client = $initialData->getClient();

        $statusChangedTemplate = new MessageFactory(new ComplaintStatusChanged());
        $reseller = $initialData->getReseller();

        foreach ($initialData->getEmployees() as $employee) {
            $employeeMessage = $statusChangedTemplate->makeMessage($employee, $reseller, $initialData);
            $result = $this->mail->send($employeeMessage);
            $this->result->addEmployeeEmailResult($result);
        }

        $clientMessage = $statusChangedTemplate->makeMessage($client, $reseller, $initialData);
        $this->result->setClientEmailResult($this->mail->send($clientMessage));
        $this->result->setClientSmsResult($this->sms->send($clientMessage));

        return $this->result;
    }
}
