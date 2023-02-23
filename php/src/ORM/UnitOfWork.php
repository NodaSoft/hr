<?php

namespace App\ORM;

class UnitOfWork
{
    public readonly EntityMetadata $entityMetadata;
    private array $values;

    public function __construct(
        public readonly ?object $entity,
        ?EntityManager $em,
        private ?UnitOfWorkState $state = UnitOfWorkState::NEW
    ){
        $this->entityMetadata = $em->getEntityMetadata($entity);
        $this->values = $this->entityMetadata->getValues();
    }

    public function isModified(): bool {
        foreach($this->entityMetadata->getValues() as $prop => $value) {
            if($this->values[$prop] !== $value) {
                return true;
            }
        }
        return false;
    }

    public function getState(): UnitOfWorkState {
        return $this->state;
    }

    public function setState(UnitOfWorkState $state): void {
        $this->state = $state;
    }
}