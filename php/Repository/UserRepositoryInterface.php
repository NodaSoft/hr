<?php

declare(strict_types=1);

namespace NodaSoft\Repository;

use NodaSoft\DTO\NewUser;
use NodaSoft\DTO\User;
use NodaSoft\Exception\UserRepositoryException;

interface UserRepositoryInterface
{
    /**
     * Добавляет пользователя в базу данных
     *
     * @param NewUser $user
     * @return int ID добавленного пользователя
     * @throws UserRepositoryException
     */
    public function add(NewUser $user): int;

    /**
     * Возвращает список пользователей старше заданного возраста
     *
     * @param int $ageFrom
     * @param int $limit
     * @return User[]
     * @throws UserRepositoryException
     */
    public function getByAge(int $ageFrom, int $limit): array;

    /**
     * Возвращает пользователя по имени
     *
     * @param string $name
     * @return User|null
     * @throws UserRepositoryException
     */
    public function getByName(string $name): ?User;
}
