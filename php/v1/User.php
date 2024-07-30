<?php

namespace Model;

use DB\DBWorker;
use PDO;

class User
{

    /**
     * Возвращает список пользователей старше заданного возраста.
     * @param  int  $ageFrom
     * @param  int  $limit
     * @return array
     */
    public static function getUsers(int $ageFrom, int $limit = 100): array
    {
        $stmt = DBWorker::getInstance()
            ->prepare("SELECT id, name, lastName, age FROM Users WHERE age > :ageFrom LIMIT :limit");
        $stmt->execute([':ageFrom' => $ageFrom, ':limit' => $limit]);
        $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);
        $users = [];
        foreach ($rows as $row) {
            $users[] = [
                'id' => $row['id'],
                'name' => $row['name'],
                'lastName' => $row['lastName'],
                'age' => $row['age'],
            ];
        }

        return $users;
    }

    /**
     * Возвращает пользователя по имени.
     * @param  string  $name
     * @return array
     */
    public static function user(string $name): array
    {
        $stmt = DBWorker::getInstance()->prepare("SELECT id, name, lastName, age FROM Users WHERE name = :name");
        $stmt->execute([':name' => $name]);
        $user_by_name = $stmt->fetch(PDO::FETCH_ASSOC);

        return [
            'id' => $user_by_name['id'],
            'name' => $user_by_name['name'],
            'lastName' => $user_by_name['lastName'],
            'age' => $user_by_name['age'],
        ];
    }

    /**
     * Добавляет пользователя в базу данных.
     * @param  string  $name
     * @param  string  $lastName
     * @param  int  $age
     * @return string
     */
    public static function add(string $name, string $lastName, int $age): string
    {
        $sth = DBWorker::getInstance()->prepare("INSERT INTO Users (name, age, lastName) VALUES (:name, :age, :lastName)");
        $sth->execute([':name' => $name, ':age' => $age, ':lastName' => $lastName]);

        return DBWorker::getInstance()->lastInsertId();
    }
}