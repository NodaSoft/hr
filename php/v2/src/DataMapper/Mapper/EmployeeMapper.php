<?php

namespace NodaSoft\DataMapper\Mapper;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\EntityInterface\Entity;

class EmployeeMapper implements Mapper
{
    /**
     * @return null|Employee
     */
    public function getById(int $id): ?Entity
    {
        // TODO: Implement getById() method.
    }

    /**
     * @param int[] $employeeIds
     * @return EmployeeCollection
     */
    public function getAllByIds(array $employeeIds): EmployeeCollection
    {
        // TODO: Implement getAllByIds() method
    }
}
