<?php

namespace Gateway;

use PDO;

class User
{
    private static PDO $instance;

    /**
     * Реализация singleton
     * @return PDO
     */
    public static function getInstance(): PDO
    {
        if (is_null(self::$instance)) {
            $dsn = 'mysql:dbname=db;host=127.0.0.1';
            $user = 'dbuser';
            $password = 'dbpass';
            self::$instance = new PDO($dsn, $user, $password);
        }

        return self::$instance;
    }

    /**
     * Возвращает список пользователей старше заданного возраста.
     * @param int $age
     * @param int $limit
     * @return array
     */
    public static function getOlderThan(int $age, int $limit): array
    {
        $stmt = self::executeSql(
            'SELECT id, name, lastName, `from`, age, settings FROM Users WHERE age > :age LIMIT :limit',
            compact('age', 'limit')
        );

        $users = $stmt->fetchAll(PDO::FETCH_ASSOC);

        foreach ($users as $i => $user) {
            $settings = json_decode($user['settings']);
            $users[$i]['key'] = $settings->key;
        }

        return $users;
    }

    /**
     * Возвращает пользователя по имени.
     * @param string $name
     * @return array
     */
    public static function findByName(string $name): array
    {
        $stmt = self::executeSql(
            'SELECT id, name, lastName, `from`, age FROM Users WHERE name = :name',
            compact('name')
        );

        return $stmt->fetch(PDO::FETCH_ASSOC);
    }

    /**
     * Добавляет пользователя в базу данных.
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @return string
     */
    public static function create(string $name, string $lastName, int $age): string
    {
        self::executeSql(
            'INSERT INTO Users (name, lastName, age) VALUES (:name, :lastName, :age)',
            compact('name', 'lastName', 'age')
        );

        return self::getInstance()->lastInsertId();
    }

    /**
     * Добавляет несколько пользователей в базу данных.
     * @param \Closure $callback
     * @return array
     * @throws \Exception
     */
    public static function transaction(\Closure $callback): array
    {
        \Gateway\User::getInstance()->beginTransaction();

        try {
            $result = $callback();

            \Gateway\User::getInstance()->commit();
        } catch (\Exception $e) {
            \Gateway\User::getInstance()->rollBack();
            throw $e;
        }

        return $result;
    }

    /**
     * Выполняет sql
     * @param string $sql
     * @param array $params
     * @return false|\PDOStatement
     */
    private static function executeSql(string $sql, array $params = []): bool|\PDOStatement
    {
        $stmt = self::getInstance()->prepare($sql);
        $stmt->execute($params);

        return $stmt;
    }
}