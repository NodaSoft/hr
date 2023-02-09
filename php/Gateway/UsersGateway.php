<?php

namespace Gateways;

use System\Database\PdoProvider;

class UsersGateway
{
    const LIMIT = 10;

    /**
     * Возвращает список пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    public static function getUsersOlderThan(int $ageFrom): array
    {
        $stmt = PdoProvider::getInstance()->prepare(
            "SELECT id, name, lastName, \"from\", age, settings->'key' FROM Users WHERE age > ? LIMIT ?"
        );

        $stmt->execute([$ageFrom, UsersGateway::LIMIT]);
        return $stmt->fetchAll(\PDO::FETCH_ASSOC);
    }

    /**
     * Возвращает пользователя по имени.
     * @param string $name
     * @return array
     */
    public static function getUserByName(string $name): array
    {
        $stmt = PdoProvider::getInstance()->prepare(
            "SELECT id, name, lastName, \"from\", age, settings->'key' FROM Users WHERE name = ?"
        );
        $stmt->execute([$name]);
        return $stmt->fetch(\PDO::FETCH_ASSOC);
    }

    /**
     * Добавляет пользователя в базу данных.
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @return string
     */
    public static function addUser(string $name, string $lastName, int $age): string
    {
        $pdo = PdoProvider::getInstance();
        $stmt = $pdo->prepare(
            "INSERT INTO Users (name, age, lastName) VALUES (:name, :age, :lastName)"
        );
        $stmt->execute([":name" => $name, ":age" => $age, ":lastName" => $lastName]);

        // Никогда этим не пользовался, могу написать ошибки. Предпочитаю returning.
        return $pdo->lastInsertId();
    }
}
