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
        $reflect = new ReflectionClass($this->entityClass);
        $from = strtolower($reflect->getShortName());

        foreach($reflect->getAttributes(ORM\Entity::class) as $reflectionAttribute) {
            $attrs = $reflectionAttribute->getArguments();
            if(isset($attrs['table'])) {
                $from = $attrs['table'];
            }
        }

        return $this->getEntityManager()
            ->queryBuilder()
            ->from($this->entityClass, $from);
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

    public function findOnBy(array $criteria) {
        $builder = $this->queryBuilder();

        foreach($criteria as $param => $value) {
            $builder->where($param, $value);
        }

        return $builder->fetchOne();
    }

    public function findAll(array $criteria, int $limit = null): array {

        $builder = $this->queryBuilder();

        foreach($criteria as $param => $value) {
            $builder->where($param, $value);
        }

        if($limit) {
            $builder->limit($limit);
        }

        return $builder->fetchAll();
    }
}
