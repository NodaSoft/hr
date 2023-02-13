<?php

declare(strict_types=1);

namespace App\Service;

use App\Domain\Criteria;
use App\Domain\UserRepository;

final class UserService
{
    private const USERS_LIMIT = 10;

    public function __construct(private UserRepository $userRepo)
    {
    }

    /**
     * Возвращает пользователей старше заданного возраста.
     * 
     * @return User[]
     */
    public function getUsersFromAge(int $age): array
    {
        $criteria = new Criteria();
        $criteria->fromAge = $age;
        return $this->userRepo->getUsersByCriteria($criteria, UserService::USERS_LIMIT);
    }

    /**
     * Возвращает пользователей по списку имен.
     * 
     * @param string[] $names
     * @return User[]
     */
    public function getUserByNames(array $names): array
    {
        $criteria = new Criteria();
        $criteria->names = $names;
        return $this->userRepo->getUsersByCriteria($criteria);
    }

    /**
     * Добавляет пользователей.
     * 
     * @param User[] $users
     * @return array
     */
    public function addUsers($users): array
    {
        return $this->userRepo->addUsers($users);
    }
}