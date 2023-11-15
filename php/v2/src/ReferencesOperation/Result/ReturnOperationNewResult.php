<?php

namespace NodaSoft\ReferencesOperation\Result;

use NodaSoft\Messenger\Result;
use NodaSoft\Messenger\ResultCollection;

class ReturnOperationNewResult implements ReferencesOperationResult
{
    /** @var ResultCollection */
    private $employeeEmails;

    public function __construct()
    {
        $this->employeeEmails = new ResultCollection();
    }

    public function toArray(): array
    {
        return [
            'employeeEmails' => $this->employeeEmails->toArray(),
        ];
    }

    public function getEmployeeEmails(): ResultCollection
    {
        return $this->employeeEmails;
    }

    public function addEmployeeEmailResult(Result $result): void
    {
        $this->employeeEmails->add($result);
    }
}
