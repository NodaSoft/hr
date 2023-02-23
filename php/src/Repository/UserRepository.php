<?php

namespace App\Repository;

use App\Entity\User;
use App\ORM\EntityManager;
use App\ORM\OrderBy;

/**
 * @method User find(int $id)
 * @method User findOnBy(array $criteria, array $orderBy = [])
 * @method User[] findAll(array $criteria = [], array $orderBy = [], int $limit = null)
 */
class UserRepository extends \App\AbstractRepository
{
    public function __construct(EntityManager $entityManager)
    {
        parent::__construct($entityManager, User::class);
    }

    public function findByName(string $name): User {
        return $this->findOnBy(['name' => $name]);
    }

    /**
     * @param int $age
     * @param int|null $limit
     * @return User[]
     */
    public function getUsers(int $age, int $limit = null): array {

        $builder = $this->queryBuilder();
        $builder->where('age', '>', $age);

        if($limit) {
            $builder->limit($limit);
        }

        // $builder->orderBy('id', OrderBy::DESC);
        return $builder->fetchAll();
    }

    /**
     * @param array $names
     * @return User[]
     */
    public function getUsersByNames(array $names): array {
        return $this->queryBuilder()
            ->where('name', $names)
            ->fetchAll();
    }
}