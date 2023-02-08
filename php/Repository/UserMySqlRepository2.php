<?php

declare(strict_types=1);

namespace NodaSoft\Repository;

use NodaSoft\DTO\NewUser;
use NodaSoft\DTO\User;
use NodaSoft\Exception\UserRepositoryException;
use PDO;
use PDOException;

/**
 * Извлечение поля "key" из JSON-поля "settings" производится средствами MySQL
 */
final class UserMySqlRepository2 implements UserRepositoryInterface
{
    /**
     * @param PDO $dbh
     */
    public function __construct(
        private PDO $dbh,
    ) {
    }

    /**
     * Добавляет пользователя в базу данных
     *
     * @param NewUser $user
     * @return int ID добавленного пользователя
     * @throws UserRepositoryException
     */
    public function add(NewUser $user): int
    {
        try {
            $stmt = $this->dbh->prepare(<<<'SQL'
                INSERT INTO
                    `users` (`name`, `last_name`, `age`)
                VALUES
                    (:name, :last_name, :age)
            SQL);
            $stmt->bindValue(':name', $user->name);
            $stmt->bindValue(':last_name', $user->lastName);
            $stmt->bindValue(':age', $user->age, PDO::PARAM_INT);
            $stmt->execute();

            return (int) $this->dbh->lastInsertId();
        } catch (PDOException $exception) {
            throw new UserRepositoryException($exception->getMessage(), (int) $exception->getCode(), $exception);
        }
    }

    /**
     * Возвращает список пользователей старше заданного возраста
     *
     * @param int $ageFrom
     * @param int $limit
     * @return User[]
     * @throws UserRepositoryException
     */
    public function getByAge(int $ageFrom, int $limit): array
    {
        try {
            $stmt = $this->dbh->prepare(<<<'SQL'
                SELECT
                    `id`,
                    `name`,
                    `last_name`,
                    `age`,
                    `from`,
                    `settings`->>"$.key" AS `key`
                FROM
                    `users`
                WHERE
                    `age` > :age_from
                LIMIT
                    :limit
            SQL);
            $stmt->bindValue(':age_from', $ageFrom, PDO::PARAM_INT);
            $stmt->bindValue(':limit', $limit, PDO::PARAM_INT);
            $stmt->execute();
            $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

            $users = [];
            foreach ($rows as $row) {
                $users[] = new User(
                    id: (int) $row['id'],
                    name: $row['name'],
                    lastName: $row['last_name'],
                    age: (int) $row['age'],
                    from: $row['from'],
                    key: $row['key'],
                );
            }

            return $users;
        } catch (PDOException $exception) {
            throw new UserRepositoryException($exception->getMessage(), (int) $exception->getCode(), $exception);
        }
    }

    /**
     * Возвращает пользователя по имени
     *
     * @param string $name
     * @return User|null
     * @throws UserRepositoryException
     */
    public function getByName(string $name): ?User
    {
        try {
            $stmt = $this->dbh->prepare(<<<'SQL'
                SELECT
                    `id`,
                    `name`,
                    `last_name`,
                    `age`,
                    `from`,
                    `settings`->>"$.key" AS `key`
                FROM
                    `users`
                WHERE
                    `name` = :name
                LIMIT
                    1
            SQL);
            $stmt->bindValue(':name', $name);
            $stmt->execute();
            $row = $stmt->fetch(PDO::FETCH_ASSOC);

            return $row !== false
                ? new User(
                    id: (int) $row['id'],
                    name: $row['name'],
                    lastName: $row['last_name'],
                    age: (int) $row['age'],
                    from: $row['from'],
                    key: $row['key'],
                )
                : null;
        } catch (PDOException $exception) {
            throw new UserRepositoryException($exception->getMessage(), (int) $exception->getCode(), $exception);
        }
    }
}
