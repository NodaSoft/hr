<?php

namespace NodaSoft\Operation\Result;

use NodaSoft\Messenger;

class NotifyComplaintNewResult implements Result
{
    /** @var Messenger\ResultCollection */
    private $employeeEmails;

    public function __construct()
    {
        $this->employeeEmails = new Messenger\ResultCollection();
    }

    public function toArray(): array
    {
        return [
            'employeeEmails' => $this->employeeEmails->toArray(),
        ];
    }

    public function getEmployeeEmails(): Messenger\ResultCollection
    {
        return $this->employeeEmails;
    }

    public function addEmployeeEmailResult(Messenger\Result $result): void
    {
        $this->employeeEmails->add($result);
    }
}
