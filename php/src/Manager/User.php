<?php

namespace App\Manager;

use App\ORM\EntityManager;
use \App\Repository\UserRepository;
use \App\Entity;


class User
{
    private readonly UserRepository $repo;

    public function __construct()
    {
        $this->repo = new UserRepository(new EntityManager());
    }

    /**
     * Возвращает пользователей старше заданного возраста.
     * @return Entity\User[]
     */
    function getUsers(int $ageFrom): array {
        return $this->repo->getUsers($ageFrom);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @return Entity\User[]
     */
    public function getByNames(): array
    {
        if(!isset($_GET['names'])) {
            return [];
        }

        $names = array_filter((array)$_GET['names'], 'is_string');
        return $this->repo->getUsersByNames($names);
    }

    /**
     * Добавляет пользователей в базу данных.
     * @return int[]
     */
    public function addUsers(Entity\User ...$users): array
    {
        $em = $this->repo->getEntityManager();

        foreach ($users as $user) {
            $em->persist($user);
        }

        $ids[] = $em->flush();
        return $ids;
    }
}