<?php

namespace NodaSoft\DataMapper\Mapper;

use NodaSoft\DataMapper\Entity\Complaint;
use NodaSoft\DataMapper\EntityInterface\Entity;

class ComplaintMapper implements Mapper
{
    /**
     * @param int $id
     * @return Complaint|null
     */
    public function getById(int $id): ?Entity
    {
        // TODO: Implement getById() method.
    }
}
