<?php

namespace NodaSoft\DataMapper\Mapper;

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
     * @param int $resellerId
     * @return Employee[]
     */
    public function getAllByReseller(int $resellerId): array
    {
        // todo: Implement getAllByReseller() method.
    }
}
