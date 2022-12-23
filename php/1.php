<?php

declare(strict_types=1);

namespace Manager;

use Exception;
use Gateway\User as Gateway;

class User
{
    /**
     * Возвращает пользователей старше заданного возраста.
     * @param int $age
     * @return array
     * @throws Exception
     */
    public static function getUsers(int $age): array
    {
        return Gateway::getUsers($age);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @return array
     * @throws Exception
     */
    public static function getByNames(): array
    {
        $names = $_GET['names'] ?? null;

        if(is_array($names)){
            $names = array_filter(
                $names,
                static function (string $name) {
                    $name = filter_var($name, FILTER_SANITIZE_STRING);
                    return '' !== trim($name);
                }
            );

            foreach ($names as $name) {
                $users[] = Gateway::getUser($name);
            }
        }

        return $users ?? [];
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param array $users
     * @return array
     * @throws Exception
     */
    public static function addUsers(array $users): array
    {
        try {
            foreach ($users as $user) {
                $userId = Gateway::add($user['name'], $user['lastName'], $user['age']);

                if('' !== $userId){
                    $result[] = $userId;
                }
            }
        } catch (Exception $e) {
            throw new Exception($e->getMessage());
        }

        return $result ?? [];
    }
}