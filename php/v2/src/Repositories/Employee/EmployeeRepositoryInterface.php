<?php

namespace Nodasoft\Testapp\Repositories\Employee;


use Nodasoft\Testapp\Entities\Employee\Employee;

interface EmployeeRepositoryInterface
{
    public function getById(int $id): Employee;
}