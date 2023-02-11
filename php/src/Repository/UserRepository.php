<?php

namespace App\Repository;

use App\Entity\User;
use App\ORM\EntityManager;

/**
 * @method User find(int $id)
 * @method User findOnBy(array $criteria)
 * @method User[] findAll(array $criteria, int $limit = null)
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