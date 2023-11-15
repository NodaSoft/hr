<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\Messenger\Recipient;
use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class Reseller implements Entity, Recipient
{
    use EntityTrait\MessageRecipientEntity;

    /** @var EmployeeCollection */
    private $employees;

    public function __construct(
        int $id = null,
        string $name = null,
        string $email = null,
        int $cellphone = null,
        EmployeeCollection $employees = null
    ) {
        if ($id) $this->setId($id);
        if ($name) $this->setName($name);
        if ($email) $this->setEmail($email);
        if ($cellphone) $this->setCellphone($cellphone);
        $this->employees = $employees ?: new EmployeeCollection;
    }

    public function getEmployees(): EmployeeCollection
    {
        return $this->employees;
    }

    public function setEmployees(EmployeeCollection $employees): void
    {
        $this->employees = $employees;
    }

    public function addEmployee(Employee $employee): void
    {
        $this->employees->add($employee);
    }

    public function getFullName(): string
    {
        return $this->getName() . ' ' . $this->getId();
    }
}
