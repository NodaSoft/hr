<?php

namespace Services;

use Repositories\UserRepository;

class UserService
{
    /**
     * @var UserRepository
     */
    private UserRepository $repository;

    /**
     * UserService constructor.
     */
    public function __construct()
    {
        $this->repository = new UserRepository();
    }

    /**
     * Возвращает пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    function getUsers(int $ageFrom): array
    {
        return $this->repository->getUsersAgeFrom($ageFrom);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @return array|null
     */
    public function getByNames(): ?array
    {
        if (!empty($_GET['names'])) {
            return $this->repository->getUsersByName($_GET['names']);
        }

        return null;
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param $users
     * @return array
     * @throws \Exception
     */
    public function addUsers($users): array
    {
        return $this->repository->addUsers($users);
    }
}