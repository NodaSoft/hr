<?php

namespace NodaSoft\DataMapper\Mapper;

use NodaSoft\DataMapper\EntityInterface\Entity;

interface Mapper
{
    public function getById(int $id): ?Entity;
}
