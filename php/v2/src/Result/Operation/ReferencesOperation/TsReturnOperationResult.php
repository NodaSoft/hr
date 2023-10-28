<?php

namespace NodaSoft\Result\Operation\ReferencesOperation;

use NodaSoft\Result\Notification\NotificationResult;

class TsReturnOperationResult implements ReferencesOperationResult
{
    /** @var NotificationResult */
    private $employeeEmail;

    /** @var NotificationResult */
    private $clientEmail;

    /** @var NotificationResult */
    private $clientSms;

    public function __construct()
    {
        $this->employeeEmail = new NotificationResult();
        $this->clientEmail = new NotificationResult();
        $this->clientSms = new NotificationResult();
    }

    public function markEmployeeEmailSent(): void
    {
        $this->employeeEmail->setIsSent(true);
    }

    public function markClientEmailSent(): void
    {
        $this->clientEmail->setIsSent(true);
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
            'employeeEmail' => $this->employeeEmail->isSent(),
            'clientEmail' => $this->clientEmail->isSent(),
            'clientSms' => $this->clientSms->toArray(),
        ];
    }
}
