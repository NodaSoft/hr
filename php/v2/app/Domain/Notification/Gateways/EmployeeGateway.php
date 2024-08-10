<?php

namespace app\Domain\Notification\Gateways;

use app\Domain\Notification\Models\Employee;

class EmployeeGateway
{
    public function getEmployeeById(int $id): ?Employee
    {
        return \DB::raw("(select * from employee where id=$id)")->first();
    }

    public function getEmployeeByIdAndByType(int $id, int $type = Employee::TYPE_CONTRACTOR): ?Employee
    {
        return \DB::raw("(select * from employee where id=$id and type=$type})")->first();
    }
}