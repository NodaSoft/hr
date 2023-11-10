<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class Reseller implements Entity, MessageRecipientEntity
{
    use EntityTrait\MessageRecipientEntity;

    public function getFullName(): string
    {
        return $this->getName() . ' ' . $this->getId();
    }
}
