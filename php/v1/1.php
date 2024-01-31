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
    function getUsers(int $ageFrom): array
    {
        return \Gateway\User::getUsers($ageFrom);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @param array $names
     * @return array
     */
    public static function getByNames(array $names): array //TODO исходя из названия метода будем принимать массив
    {
        $users = [];
        foreach ($names as $name) {
            $users[] = \Gateway\User::user(trim($name));
        }

        return $users;
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param array $users
     * @return array
     */
    public function users(array $users): array
    {
        $ids = [];
        \Gateway\User::getInstance()->beginTransaction();
        foreach ($users as $user) {
            try {
                $this->validate_fields($user);
                $ids[] = \Gateway\User::add((string)$user['name'], (string)$user['lastName'], (int)$user['age']);
                \Gateway\User::getInstance()->commit();
            } catch (\Exception $e) {
                \Gateway\User::getInstance()->rollBack();
            }
        }

        return $ids;
    }

    /**
     * @param array $user
     * @return void
     * @throws \Exception
     */
    private function validate_fields(array $user): void
    {
        $fields = ['name', 'lastName', 'age'];
        foreach ($fields as $field) {
            if (!array_key_exists($field, $user)) {
                throw new \Exception('empty field ' . $field);
            }
        }
    }
}