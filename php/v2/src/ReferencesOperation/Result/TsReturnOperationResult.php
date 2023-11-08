<?php

namespace NodaSoft\ReferencesOperation\Result;

use NodaSoft\Mail\Result;
use NodaSoft\Mail\ResultCollection;
use NodaSoft\Result\Notification\NotificationResult;

class TsReturnOperationResult implements ReferencesOperationResult
{
    /** @var ResultCollection */
    private $employeeEmails;

    /** @var Result */
    private $clientEmail;

    /** @var NotificationResult */
    private $clientSms;

    public function __construct()
    {
        $this->employeeEmails = new ResultCollection();
        $this->clientSms = new NotificationResult();
    }

    public function markClientSmsSent(): void
    {
        $this->clientSms->setIsSent(true);
    }

    public function setClientSmsErrorMessage(string $message): void
    {
        $this->clientSms->setErrorMessage($message);
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

    public function getClientSms(): NotificationResult
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
}
