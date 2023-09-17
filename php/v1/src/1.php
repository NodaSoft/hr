<?php

declare(strict_types=1);

namespace Manager;

use Exception;

class User
{
    public const limit = 10;

    /**
     * Возвращает пользователей старше заданного возраста.
     */
    public function getUsers(int $ageFrom): array
    {
        return \Gateway\User::getUsers($ageFrom);
    }

    /**
     * Возвращает пользователей по списку имен.
     */
    public static function getByNames(): array
    {
        $users = [];
        foreach ($_GET['names'] as $name) {
            $users[] = \Gateway\User::user($name);
        }

        return $users;
    }

    /**
     * Добавляет пользователей в базу данных.
     */
    public function users($users): array
    {
        $ids = [];
        \Gateway\User::getInstance()->beginTransaction();
        foreach ($users as $user) {
            try {
                \Gateway\User::add($user['name'], $user['lastName'], $user['age']);
                \Gateway\User::getInstance()->commit();
                $ids[] = \Gateway\User::getInstance()->lastInsertId();
            } catch (Exception $e) {
                \Gateway\User::getInstance()->rollBack();
            }
        }

        return $ids;
    }
}
