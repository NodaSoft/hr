<?php

namespace Managers;

use Gateways\UsersGateway;
use System\Database\PdoProvider;

class UsersManager
{
    /**
     * Возвращает пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    public function getUsersOlderThan(int $ageFrom): array
    {
        return UsersGateway::getUsersOlderThan($ageFrom);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @param array $user_names
     * @return array
     */
    public static function getUsersByNames(array $user_names): array
    {
        $users = [];

        foreach ($user_names as $name) {
            $users[] = UsersGateway::getUserByName((string)$name);
        }

        return $users;
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param array $users
     * @return array
     */
    public function addUsers(array $users): array
    {
        $ids = [];

        $pdo = PdoProvider::getInstance();
        $pdo->beginTransaction();

        try {
            foreach ($users as $user) {
                // По хорошему бы сделать класс User, а то во-первых поля может не оказаться в массиве
                // можно использовать ?? но во-вторых непонятно как воспринимать пустые строки
                $ids[] = UsersGateway::addUser($user["name"], $user["lastName"], $user["age"]);
            }

            $pdo->commit();
        } catch (\Exception $e) {
            // Надо что-то делать с этим эксепшеном. Либо как-то выталкивать выше код ошибки и
            // сообщение, либо не перехватывать его здесь совсем, а предоставить это внешней обработке
            // Я предпочитаю второй вариант.
            $pdo->rollBack();
            return [];
        }

        return $ids;
    }
}
