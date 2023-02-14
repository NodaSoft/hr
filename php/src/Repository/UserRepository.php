<?php

namespace App\Repository;


use App\Connection\Repository\AbstractRepository;

/**
 *
 */
class UserRepository extends AbstractRepository
{
    /**
     * Возвращает список пользователей старше заданного возраста.
     *
     * @param int $age
     * @param int $limit
     * @return array
     * @throws \Exception
     */
    public function listByAge(int $age, int $limit): array
    {
        $statement = $this->getConnection()->prepare('SELECT id, name, lastName, `from`, age, settings FROM `Users` WHERE age > :age LIMIT :limit');

        $statement->bindParam(':age', $age, \PDO::PARAM_INT);
        $statement->bindParam(':limit', $limit, \PDO::PARAM_INT);
        $statement->execute();

        return $statement->fetchAll();
    }

    /**
     * Возвращает пользователя по имени.
     *
     * @param string $name
     * @return array
     * @throws \Exception
     */
    public function getByName(string $name)
    {
        $statement = $this->getConnection()->prepare('SELECT id, name, lastName, `from`, age FROM `Users` WHERE name = :name');

        $statement->bindParam(':name', $name);
        $statement->execute();

        return $statement->fetchAll();
    }

    /**
     * Добавляет пользователя в базу данных.
     *
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @return string|false
     * @throws \Exception
     */
    public function add(string $name, string $lastName, int $age)
    {
        $statement = $this->getConnection()->prepare("INSERT INTO `Users` (name, lastName, age) VALUES (:name, :lastName, :age)");

        $statement->bindParam(':name', $name);
        $statement->bindParam(':lastName', $lastName);
        $statement->bindParam(':age', $age, \PDO::PARAM_INT);

        $statement->execute();

        return (int)$this->getConnection()->lastInsertId();
    }
}