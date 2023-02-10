<?php

namespace Repositories;

use PDO;
use db\DBInstance;

class UserRepository
{
    private PDO $DBInstance;

    const GET_USERS_AGE_FROM_LIMIT = 10;

    /**
     * UserRepository constructor.
     */
    public function __construct()
    {
        $this->DBInstance = DBInstance::getInstance();
    }

    /**
     * Возвращает список пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    public function getUsersAgeFrom(int $ageFrom): array
    {
        $stmt = $this->DBInstance->prepare("SELECT id, name, lastName, from, age, settings FROM Users WHERE age > :age LIMIT :limit");
        $stmt->execute(['age' => $ageFrom, 'limit' => self::GET_USERS_AGE_FROM_LIMIT]);

        $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

        $users = [];
        foreach ($rows as $row) {
            $settings = json_decode($row['settings']);
            $users[] = [
                'id' => $row['id'],
                'name' => $row['name'],
                'lastName' => $row['lastName'],
                'from' => $row['from'],
                'age' => $row['age'],
                'key' => $settings['key'],
            ];
        }

        return $users;
    }

    /**
     * Возвращает пользователя по имени.
     * @param string $name
     * @return array|null
     */
    public function getUserByName(string $name): ?array
    {
        $stmt = $this->DBInstance->prepare("SELECT id, name, lastName, from, age, settings FROM Users WHERE name = :name");
        $stmt->execute(['name' => $name]);
        $userByName = $stmt->fetch(PDO::FETCH_ASSOC);

        if (!empty($userByName)) {
            return [
                'id' => $userByName['id'],
                'name' => $userByName['name'],
                'lastName' => $userByName['lastName'],
                'from' => $userByName['from'],
                'age' => $userByName['age'],
            ];
        }

        return null;
    }

    /**
     * Возвращает пользователей по имени.
     * @param array $names
     * @return array|null
     */
    public function getUsersByName(array $names): ?array
    {
        $stmt = $this->DBInstance->prepare("SELECT id, name, lastName, from, age, settings FROM Users WHERE name IN (:name)");
        $stmt->execute(['name' => explode(',', $names)]);

        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    /**
     * Добавляет пользователя в базу данных.
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @return int
     */
    public function addUser(string $name, string $lastName, int $age): int
    {
        $sth = $this->DBInstance->prepare("INSERT INTO Users (name, age, lastName) VALUES (:name, :age, :lastName)");
        $sth->execute([':name' => $name, ':age' => $age, ':lastName' => $lastName]);

        return $this->DBInstance->lastInsertId();
    }

    /**
     * @param array $users
     * @return array
     */
    public function addUsers(array $users)
    {
        $this->DBInstance->beginTransaction();
        $ids = [];

        foreach ($users as $user) {
            if (isset($user['name'], $user['lastName'], $user['age'])) {
                try {
                    $ids[] = $this->addUser($user['name'], $user['lastName'], $user['age']);
                } catch (\Exception $exception){
                    $this->DBInstance->rollBack();
                }
            }
        }

        $this->DBInstance->commit();

        return $ids;
    }
}