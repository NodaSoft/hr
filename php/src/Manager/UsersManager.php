<?php

namespace Manager;

use Db\Db;
use Gateway\UsersGateway;

class UsersManager
{
    public const LIMIT = 10;

    /**
     * Возвращает пользователей старше заданного возраста.
     *
     * @param int $ageFrom
     *
     * @return array
     */
    public function getByAge(int $ageFrom): array
    {
        return (new UsersGateway())->byAge($ageFrom)->all(self::LIMIT);
    }

    /**
     * Возвращает пользователей по списку имен.
     *
     * @param array $names
     *
     * @return array
     */
    public function getByNames(array $names): array
    {
        return (new UsersGateway())->byName($names)->all();
    }

    /**
     * Добавляет пользователей в базу данных.
     *
     * @param array $users
     *
     * @return array
     * @throws \Throwable
     */
    public function insertUsers(array $users): array
    {
        $ids = [];
        $pdo = Db::i()->pdo();
        $pdo->beginTransaction();
        foreach ($users as $user) {
            try {
                $ids[] = (new UsersGateway())->add($user['name'], $user['lastName'], (int)$user['age']);
            } catch (\Throwable $throwable) {
                $pdo->rollBack();
                throw $throwable;
            }
        }
        $pdo->commit();

        return $ids;
    }
}
