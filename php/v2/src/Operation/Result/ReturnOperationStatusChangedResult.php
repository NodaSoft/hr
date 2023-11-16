<?php

namespace NodaSoft\Operation\Result;

use NodaSoft\Messenger;

class ReturnOperationStatusChangedResult implements Result
{
    /** @var Messenger\ResultCollection */
    private $employeeEmails;

    /** @var ?Messenger\Result */
    private $clientEmail = null;

    /** @var ?Messenger\Result */
    private $clientSms = null;

    public function __construct()
    {
        $this->employeeEmails = new Messenger\ResultCollection();
    }

    /**
     * @return array<string, mixed>
     */
    public function toArray(): array
    {
        return [
            'employeeEmails' => $this->employeeEmails->toArray(),
            'clientEmail' => $this->clientEmail ? $this->clientEmail->toArray() : null,
            'clientSms' => $this->clientSms ? $this->clientSms->toArray() : null,
        ];
    }

    public function getEmployeeEmails(): Messenger\ResultCollection
    {
        return $this->employeeEmails;
    }

    public function getClientEmail(): ?Messenger\Result
    {
        return $this->clientEmail;
    }

    public function getClientSms(): ?Messenger\Result
    {
        return $this->clientSms;
    }

    public function addEmployeeEmailResult(Messenger\Result $result): void
    {
        $this->employeeEmails->add($result);
    }

    public function setClientEmailResult(Messenger\Result $result): void
    {
        $this->clientEmail = $result;
    }

    public function setClientSmsResult(Messenger\Result $result): void
    {
        $this->clientSms = $result;
    }
}
