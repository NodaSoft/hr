<?php

declare(strict_types=1);

namespace Repository;


use Exception;
use Model\User;


class UserRepository
{
    const int limit = 10;

    /**
     * Возвращает пользователей старше заданного возраста.
     * @param  int  $ageFrom
     * @return array
     */
    function getUsers(int $ageFrom): array
    {
        return User::getUsers($ageFrom);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @return array
     */
    public static function getByNames(string ...$names): array
    {
        $users = [];
        foreach ($names as $name) {
            $users[] = User::user($name);
        }

        return $users;
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param  array  $users
     * @return array
     * @throws Exception
     */
    public function users(array $users): array
    {
        $ids = [];
        foreach ($users as $user) {

            if (!isset($user['name'], $user['lastName'], $user['age'])) {
                continue;
            }

            try {
                $ids[] = User::add($user['name'], $user['lastName'], $user['age']);
            } catch (Exception $e) {
                throw new Exception("Failed to add user '{$user['name']}'.", 500, $e);
            }
        }

        return $ids;
    }
}