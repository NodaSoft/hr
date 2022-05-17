<?php

namespace Manager;

class User
{
    const limit = 10;

    /**
     * Возвращает пользователей старше заданного возраста.
     * @param int $age
     * @return array
     */
    public static function getOlderThan(int $age): array
    {
        return \Gateway\User::getOlderThan($age, self::limit);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @param array $names
     * @return array
     */
    public static function getByNames(array $names): array
    {
        $users = [];

        foreach ($names as $name) {
            $users[] = \Gateway\User::findByName($name);
        }

        return $users;
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param $users
     * @return array
     * @throws \Exception
     */
    public static function create($users): array
    {
        return \Gateway\User::transaction(function () use ($users) {
            $ids = [];

            foreach ($users as $user) {
                $ids[] = \Gateway\User::create($user['name'], $user['lastName'], $user['age']);
            }

            return $ids;
        });
    }
}