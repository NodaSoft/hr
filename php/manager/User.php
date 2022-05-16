<?php

namespace Manager;

class User
{

    const LIMIT = 10;

    /**
     * Возвращает пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    public static function getUsersByAgeFrom(int $ageFrom): array
    {
        return \Gateway\User::getUsers($ageFrom);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @param array $names
     * @return array
     */
    public static function getUsersByName(array $names): array
    {
        $users = [];
        foreach ($names as $name) {
            $users[] = \Gateway\User::user($name);
        }

        return $users;
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param $users
     * @return array
     */
    public static function addUsers($users): array
    {
        $ids = [];
        try {
            \Gateway\User::getInstance()->beginTransaction();
            foreach ($users as $user) {
                $ids[] = \Gateway\User::add($user['name'], $user['lastName'], $user['age']);
            }
            \Gateway\User::getInstance()->commit();
        } catch (\Exception $e) {
            \Gateway\User::getInstance()->rollBack();
            $ids = [];
        }

        return $ids;
    }

}

