<?php

namespace NW\WebService\References\Operations\Notification\Contracts;

use NW\WebService\References\Operations\Notification\Employee;

interface EmployeeServiceContract
{
    public function getById(int $id): ?Employee;
}
