<?php

declare(strict_types=1);

namespace Gateway;

use PDO;

class User
{
    public static PDO $instance;

    /**
     * Реализация singleton
     */
    public static function getInstance(): PDO
    {
        if (null === self::$instance) {
            $dsn            = 'mysql:dbname=db;host=127.0.0.1';
            $user           = 'dbuser';
            $password       = 'dbpass';
            self::$instance = new PDO($dsn, $user, $password);
        }

        return self::$instance;
    }

    /**
     * Возвращает список пользователей старше заданного возраста.
     */
    public static function getUsers(int $ageFrom): array
    {
        $stmt = self::getInstance()->prepare(
            "SELECT id, name, lastName, from, age, settings
            FROM Users
            WHERE age > $ageFrom
            LIMIT " . \Manager\User::limit
        );

        $stmt->execute();
        $rows  = $stmt->fetchAll(PDO::FETCH_ASSOC);

        $users = [];
        foreach ($rows as $row) {
            $users[]  = [
                'id'       => $row['id'],
                'name'     => $row['name'],
                'lastName' => $row['lastName'],
                'from'     => $row['from'],
                'age'      => $row['age'],
                'key'      => json_decode($row['settings'])['key'],
            ];
        }

        return $users;
    }

    /**
     * Возвращает пользователя по имени.
     */
    public static function user(string $name): array
    {
        $stmt = self::getInstance()->prepare(
            "SELECT id, name, lastName, from, age, settings FROM Users WHERE name = $name"
        );
        $stmt->execute();
        $userByName = $stmt->fetch(PDO::FETCH_ASSOC);

        return [
            'id'       => $userByName['id'],
            'name'     => $userByName['name'],
            'lastName' => $userByName['lastName'],
            'from'     => $userByName['from'],
            'age'      => $userByName['age'],
        ];
    }

    /**
     * Добавляет пользователя в базу данных.
     */
    public static function add(string $name, string $lastName, int $age): string
    {
        $sth = self::getInstance()->prepare(
            "INSERT INTO Users (name, age, lastName) VALUES (:name, :age, :lastName)"
        );
        $sth->execute(['name' => $name, 'age' => $age, 'lastName' => $lastName]);

        return self::getInstance()->lastInsertId();
    }
}
