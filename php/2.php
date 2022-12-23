<?php

declare(strict_types=1);

namespace Gateway;

use Exception;
use Manager\User as Manager;
use PDO;
use PDOException;

class User
{
    /**
     * @var PDO
     */
    public static $instance = null;

    /**
     * @var string
     */
    private static $dsn = 'mysql:host=localhost;dbname=db;';

    /**
     * @var string
     */
    private static $db_user = 'dbuser';

    /**
     * @var string
     */
    private static $db_pass = 'dbpass';

    /**
     * Реализация singleton
     * @return PDO
     * @throws Exception
     */
    public static function getInstance(): PDO
    {
        if (!self::$instance) {
            try {
                self::$instance = new PDO(self::$dsn, self::$db_user, self::$db_pass);
                self::$instance->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            } catch (PDOException $e) {
                throw new Exception($e->getMessage());
            }
        }

        return self::$instance;
    }

    /**
     * Возвращает список пользователей старше заданного возраста.
     * @param int $age
     * @return array
     * @throws Exception
     */
    public static function getUsers(int $age): array
    {
        try {
            $stmt = self::getInstance()->prepare("SELECT id, name, lastName, from, age, settings FROM Users WHERE age > :age LIMIT :limit");
            $stmt->bindValue(':age', $age, PDO::PARAM_INT);
            $stmt->bindValue(':limit', Manager::limit, PDO::PARAM_INT);
            $stmt->execute();
            $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

            foreach ($rows as $row) {
                $settings = $row['settings'] ?? null;

                if(null !== $settings){
                    $settings = json_decode($settings);
                }

                $result[] = [
                    'id' => $row['id'],
                    'name' => $row['name'],
                    'lastName' => $row['lastName'],
                    'from' => $row['from'],
                    'age' => $row['age'],
                    'key' => $settings['key'] ?? null,
                ];
            }
        } catch (Exception $e) {
            throw new Exception($e->getMessage());
        }

        return $result ?? [];
    }

    /**
     * Возвращает пользователя по имени.
     * @param string $name
     * @return array
     * @throws Exception
     */
    public static function getUser(string $name): array
    {
        try {
            $stmt = self::getInstance()->prepare("SELECT id, name, lastName, from, age FROM Users WHERE name = :name");
            $stmt->bindValue(':name', $name);
            $stmt->execute();

            if($user = $stmt->fetch(PDO::FETCH_ASSOC)){
                $result = [
                    'id' => $user['id'],
                    'name' => $user['name'],
                    'lastName' => $user['lastName'],
                    'from' => $user['from'],
                    'age' => $user['age']
                ];
            }
        } catch (Exception $e) {
            throw new Exception($e->getMessage());
        }

        return $result ?? [];
    }

    /**
     * Добавляет пользователя в базу данных.
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @return string
     * @throws Exception
     */
    public static function add(string $name, string $lastName, int $age): string
    {
        $dbh = self::getInstance();

        if ($dbh->beginTransaction()){
            try {
                $stmt = $dbh->prepare("INSERT INTO Users (name, lastName, age) VALUES (:name, :lastName, :age)");
                $stmt->bindValue(':name', $name);
                $stmt->bindValue(':lastName', $lastName);
                $stmt->bindValue(':age', $age, PDO::PARAM_INT);
                $stmt->execute();
                
                $dbh->commit();
                $id = $dbh->lastInsertId();
            } catch (Exception $e) {
                if ($dbh->inTransaction()){
                    $dbh->rollBack();
                }

                throw new Exception($e->getMessage());
            }
        }

        return $id ?? '';
    }
}