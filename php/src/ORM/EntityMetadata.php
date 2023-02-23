<?php

namespace App\ORM;

use App\ORM\Column;
use App\ORM\ColumnType;
use App\ORM\ID;
use ReflectionException;
use ReflectionObject;
use ReflectionProperty;

class EntityMetadata extends EntityClassMetadata
{
    public function __construct(private readonly object $entity)
    {
        parent::__construct(new ReflectionObject($entity));
    }

    public function getValues(): array
    {
        $values = [];

        foreach($this->reflect->getProperties() as $prop) {
            if(!($attrs = $prop->getAttributes(Column::class))) {
                continue;
            }

            /** @var Column $attr */
            $attr = $attrs[0]->newInstance();

            if($prop->isInitialized($this->entity)) {
                $values[$prop->getName()] = $prop->getValue($this->entity);
            } else {
                $values[$prop->getName()] = $prop->getDefaultValue();
            }

            if($attr->type == ColumnType::JSON) {
                $values[$prop->getName()] = json_encode($values[$prop->getName()]);
            }
        }

        return $values;
    }

    public function getIdProperty(): ?ReflectionProperty
    {
        foreach($this->reflect->getProperties() as $prop){
            if($prop->getAttributes(ID::class)){
                return $prop;
            }
        }

        return null;
    }

    public function getIdPropertyValue(): int | string | null
    {
        foreach($this->reflect->getProperties() as $prop){
            if($prop->getAttributes(ID::class)){
                return $prop->getValue($this->entity);
            }
        }

        return null;
    }

    public function getEntityObjectId(): string {
        return spl_object_id($this->entity);
    }
}