<?php

namespace App;

use ReflectionClass;

class EntityClassMetadata
{
    public readonly string $from;
    public readonly string $entityClass;
    protected ?ORM\Entity $attrEntity = null;

    public function __construct(protected readonly ReflectionClass $reflect)
    {
        $this->entityClass = $reflect->getName();

        if($reflectionAttribute = $this->reflect->getAttributes(ORM\Entity::class)) {
            $this->attrEntity = $reflectionAttribute[0]->newInstance();
        }

        if($this->attrEntity) {
            $this->from = $this->attrEntity->table ?? strtolower($this->reflect->getShortName());
        }
    }
}