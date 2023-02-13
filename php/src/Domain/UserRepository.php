<?php

declare(strict_types=1);

namespace App\Domain;

interface UserRepository
{
    /**
     * Возвращает список пользователей по заданному критерию.
     * 
     * @return User[]
     */
    public function getUsersByCriteria(Criteria $criteria, ?int $limit = null): array;

    /**
     * Добавляет пользователей в базу данных.
     * 
     * @param User[] $users
     * @return User[]
     */
    public function addUsers(array $users): array;
}
