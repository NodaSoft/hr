<?php

namespace Src\Employee\Domain\Repository;

use Src\Employee\Domain\Entity\Employee;

interface EmployeeRepositoryInterface
{
    public function getById(int $employeeId): Employee;
}