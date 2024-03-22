<?php

namespace Src\Employee\Infrastructure\Repository;

use Src\Employee\Domain\Repository\EmployeeRepositoryInterface;
use Src\Employee\Domain\Entity\Employee;

class EmployeeRepository implements EmployeeRepositoryInterface
{

    public function getById(int $employeeId): Employee
    {
        // TODO: Implement getById() method.
    }
}
