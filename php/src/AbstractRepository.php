<?php

namespace App;

use App\ORM\EntityManager;
use ReflectionClass;

abstract class AbstractRepository
{
    public function __construct(
        protected readonly EntityManager $em,
        private readonly string $entityClass
    ){
    }

    public function queryBuilder(): ORM\QueryBuilder
    {
        $em = $this->getEntityManager();
        $meta = $em->getEntityClassMetadata($this->entityClass);
        return $em->queryBuilder()->from($this->entityClass, $meta->from);
    }

    /**
     * @return EntityManager
     */
    public function getEntityManager(): EntityManager
    {
        return $this->em;
    }

    public function find(int $id) {
        return $this->queryBuilder()->where('id', $id)->fetchOne();
    }

    public function findOnBy(array $criteria, array $orderBy = []) {
        $builder = $this->queryBuilder();

        foreach($criteria as $param => $value) {
            $builder->where($param, $value);
        }

        foreach($orderBy as $col => $order) {
            $builder->orderBy($col, $order);
        }

        return $builder->fetchOne();
    }

    public function findAll(array $criteria, array $orderBy = [], int $limit = null): array {

        $builder = $this->queryBuilder();

        foreach($criteria as $param => $value) {
            $builder->where($param, $value);
        }

        foreach($orderBy as $col => $order) {
            $builder->orderBy($col, $order);
        }

        if($limit) {
            $builder->limit($limit);
        }

        return $builder->fetchAll();
    }
}
