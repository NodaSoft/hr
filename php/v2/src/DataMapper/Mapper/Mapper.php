<?php

namespace NodaSoft\DataMapper\Mapper;

use NodaSoft\DataMapper\Entity\Entity;

interface Mapper
{
    public function getById(int $id): ?Entity;
}
