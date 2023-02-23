<?php

namespace App;

use App\ORM\EntityManager;
use ReflectionClass;

abstract class AbstractRepository
{
    public function __construct(
        protected readonly EntityManager $em,
        private readonly string $entityName
    ){
    }

    public function queryBuilder(): ORM\QueryBuilder
    {
        return $this->getEntityManager()->queryBuilder()->from($this->entityName);
    }

    /**
     * @return EntityManager
     */
    public function getEntityManager(): EntityManager
    {
        return $this->em;
    }

    public function find(int $id)
    {
        return $this->em->find($this->entityName, $id);
    }

    public function findOnBy(array $criteria, array $orderBy = [])
    {
        return $this->em->findOnBy($this->entityName, ...func_get_args());
    }

    public function findAll(array $criteria = [], array $orderBy = [], int $limit = null): array
    {
        return $this->em->findAll($this->entityName, ...func_get_args());
    }
}
