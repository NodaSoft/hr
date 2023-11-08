<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\EmailEntity;
use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class Employee implements Entity, EmailEntity
{
    use EntityTrait\EmailEntity;

    public function getFullName(): string
    {
        return $this->getName() . ' ' . $this->getId();
    }
}
