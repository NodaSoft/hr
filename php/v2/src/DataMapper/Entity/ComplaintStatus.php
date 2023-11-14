<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class ComplaintStatus implements Entity
{
    use EntityTrait\Entity;

    public function __construct(
        int $id = null,
        string $name = null
    ) {
        if ($id) $this->setId($id);
        if ($name) $this->setName($name);
    }
}
