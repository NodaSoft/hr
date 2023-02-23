<?php

namespace App\ORM;

use App\Exception\EntityClassMetadataException;
use ReflectionClass;

class EntityClassMetadata
{
    public readonly string $from;
    public readonly string $entityName;
    protected ?Entity $attrEntity = null;

    /**
     * @throws EntityClassMetadataException
     */
    public function __construct(protected readonly ReflectionClass $reflect)
    {
        if(!$reflect->getAttributes(Entity::class)) {
            throw new EntityClassMetadataException(sprintf('Class %s is not entity!', $reflect->getName()));
        }

        $this->entityName = $reflect->getName();

        if($reflectionAttribute = $this->reflect->getAttributes(Entity::class)) {
            $this->attrEntity = $reflectionAttribute[0]->newInstance();
        }

        if($this->attrEntity) {
            $this->from = $this->attrEntity->table ?? strtolower($this->reflect->getShortName());
        }
    }
}