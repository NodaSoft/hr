<?php

namespace App\ORM;

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

        foreach($this->getProps() as $prop) {
            $attrs = $prop->getAttributes(Column::class);

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

    /**
     * @return \ReflectionProperty[]
     */
    public function getProps(): array
    {
        return array_filter($this->reflect->getProperties(), function($prop){
            return count($prop->getAttributes(Column::class)) > 0;
        });
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