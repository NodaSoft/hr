<?php

namespace  Nodasoft\Testapp\Repositories\Employee;

use Exception;
use Nodasoft\Testapp\Entities\Employee\Employee;
use Nodasoft\Testapp\Entities\Employee\EmployeeMockData;
use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Traits\CanGetByKey;

class EmployeeRepository implements EmployeeRepositoryInterface
{
    use CanGetByKey;

    /**
     * @throws Exception
     */
    public function getById(int $id): Employee
    {
        $record = $this->getByKeyOrThrow(EmployeeMockData::get(), $id);

        return new Employee(
            $record['id'],
            $record['name'],
            $record['email'],
            ContactorType::tryFrom($record['type'])
        );
    }
}