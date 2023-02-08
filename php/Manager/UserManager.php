<?php

declare(strict_types=1);

namespace NodaSoft\Manager;

use NodaSoft\DTO\NewUser;
use NodaSoft\DTO\User;
use NodaSoft\Exception\UserRepositoryException;
use NodaSoft\Repository\UserRepositoryInterface;

final class UserManager
{
    private const LIMIT = 10;

    /**
     * @param UserRepositoryInterface $userRepository
     */
    public function __construct(
        private UserRepositoryInterface $userRepository,
    ) {
    }

    /**
     * Добавляет пользователей в базу данных
     *
     * @param NewUser[] $users
     * @return int[] Массив ID добавленных пользователей
     * @throws UserRepositoryException
     */
    public function add(array $users): array
    {
        $ids = [];
        foreach ($users as $user) {
            $ids[] = $this->userRepository->add($user);
        }

        return $ids;
    }

    /**
     * Возвращает пользователей старше заданного возраста
     *
     * @param int $ageFrom
     * @return User[]
     * @throws UserRepositoryException
     */
    public function getByAge(int $ageFrom): array
    {
        return $this->userRepository->getByAge($ageFrom, self::LIMIT);
    }

    /**
     * Возвращает пользователей по списку имён
     *
     * @param string[] $names Массив имён
     * @return User[]
     * @throws UserRepositoryException
     */
    public function getByNames(array $names): array
    {
        $users = [];
        foreach ($names as $name) {
            $user = $this->userRepository->getByName($name);
            if (isset($user)) {
                $users[] = $user;
            }
        }

        return $users;
    }
}
