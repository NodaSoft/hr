<?php

namespace Nodasoft\Testapp\Entities\Employee;

use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Interfaces\MockDataInterface;

class EmployeeMockData implements MockDataInterface
{
    public static function get(): array
    {
        return [
            [
                'id' => 1,
                'name' => 'employee 1',
                'email' => 'employee1email@noda.soft',
                'type' => ContactorType::TYPE_EMPLOYEE->value,
            ],
            [
                'id' => 2,
                'name' => 'employee 2',
                'email' => 'employee2email@noda.soft',
                'type' => ContactorType::TYPE_EMPLOYEE->value,
            ],
            [
                'id' => 3,
                'name' => 'employee 3',
                'email' => 'employee3email@noda.soft',
                'type' => ContactorType::TYPE_EMPLOYEE->value,
            ],
        ];
    }
}