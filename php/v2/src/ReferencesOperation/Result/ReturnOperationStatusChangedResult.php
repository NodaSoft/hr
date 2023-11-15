<?php

namespace NodaSoft\ReferencesOperation\Result;

use NodaSoft\Messenger\Result;
use NodaSoft\Messenger\ResultCollection;

class ReturnOperationStatusChangedResult implements ReferencesOperationResult
{
    /** @var ResultCollection */
    private $employeeEmails;

    /** @var Result */
    private $clientEmail;

    /** @var ?Result */
    private $clientSms;

    public function __construct()
    {
        $this->employeeEmails = new ResultCollection();
    }

    public function toArray(): array
    {
        return [
            'employeeEmails' => $this->employeeEmails->toArray(),
            'clientEmail' => $this->clientEmail->toArray(),
            'clientSms' => $this->clientSms->toArray(),
        ];
    }

    public function getEmployeeEmails(): ResultCollection
    {
        return $this->employeeEmails;
    }

    public function getClientEmail(): Result
    {
        return $this->clientEmail;
    }

    public function getClientSms(): Result
    {
        return $this->clientSms;
    }

    public function addEmployeeEmailResult(Result $result): void
    {
        $this->employeeEmails->add($result);
    }

    public function setClientEmailResult(Result $result): void
    {
        $this->clientEmail = $result;
    }

    public function setClientSmsResult(Result $result): void
    {
        $this->clientSms = $result;
    }
}
