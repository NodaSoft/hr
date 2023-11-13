<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class Employee implements Entity, MessageRecipientEntity
{
    use EntityTrait\MessageRecipientEntity;

    public function __construct(
        int $id = null,
        string $name = null,
        string $email = null,
        int $cellphone = null
    ) {
        if ($id) $this->setId($id);
        if ($name) $this->setName($name);
        if ($email) $this->setEmail($email);
        if ($cellphone) $this->setCellphone($cellphone);
    }

    public function getFullName(): string
    {
        return $this->getName() . ' ' . $this->getId();
    }
}
