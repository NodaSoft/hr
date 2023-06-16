<?php

declare(strict_types=1);

namespace App\Repository;

use App\Model\UserModel;
use App\Service\PdoService;
use Exception;
use PDO;
use RuntimeException;

class UserRepository extends PdoService
{
    /**
     * Возвращает список пользователей старше заданного возраста.
     * @param int $ageFrom
     * @param int|null $limit
     * @param int|null $offset
     * @return array
     */
    public function getUserAgeFrom(int $ageFrom, int $limit = null, int $offset = null): array
    {
        return $this->fetchAll(
            $this->getSql('userAgeFrom'),
            PDO::FETCH_ASSOC,
            [':age' => $ageFrom],
            $limit,
            $offset
        );
    }

    /**
     * Возвращает пользователя по имени.
     * Имя не уникально, по этому нужно добавить лимит
     * @param string $name
     * @param int|null $limit
     * @param int|null $offset
     * @return array|null
     */
    public function getUserByName(string $name, int $limit = null, int $offset = null): ?array
    {
        return $this->fetchAll(
            $this->getSql('userByName'),
            PDO::FETCH_ASSOC,
            [':name' => $name],
            $limit,
            $offset
        );
    }

    /**
     * Возвращает пользователя по логину.
     * @param string $login
     * @return array|null
     */
    public function getUserByLogin(string $login): ?array
    {
        return $this->fetch(
            $this->getSql('userByLogin'),
            PDO::FETCH_ASSOC,
            [':login' => $login]
        );
    }

    /**
     * Возвращает пользователей по списку имен.
     * @param array $names
     * @return array
     */
    public function listUsersByName(array $names): array
    {
        $in = str_repeat('?,', count($names) - 1) . '?';

        $sth = $this
            ->connect()
            ->prepare($this->getSql('listUsersByName', $in));
        $sth->execute($names);

        $res = $sth->fetchAll(PDO::FETCH_ASSOC);
        return $res ?: [];
    }

    /**
     * Возвращает пользователей по списку логинов.
     * @param array $logins
     * @return array
     */
    public function listUsersByLogin(array $logins): array
    {
        $in = str_repeat('?,', count($logins) - 1) . '?';

        $sth = $this
            ->connect()
            ->prepare($this->getSql('listUsersByLogin', $in));
        $sth->execute($logins);

        $res = $sth->fetchAll(PDO::FETCH_ASSOC);
        return $res ?: [];
    }

    /**
     * Добавляет пользователя в базу данных.
     * @param UserModel $user
     * @return int
     */
    public function addUser(UserModel $user): int
    {
        if ($this->getUserByLogin($user->getLogin()) !== null) {
            throw new RuntimeException('User already registered');
        }

        $settings = null;
        if ($user->getSettings() !== null) {
            $settings = json_encode($user->getSettings(), JSON_THROW_ON_ERROR);
        }

        $connect = $this->connect();
        $sth = $connect
            ->prepare($this->getSql('addUser'));
        $sth->execute(
            [
                ':login' => $user->getLogin(),
                ':name' => $user->getName(),
                ':lastName' => $user->getLastName(),
                ':age' => $user->getAge(),
                ':from' => $user->getFrom(),
                ':settings' => $settings,
            ]
        );

        return (int)$connect->lastInsertId();
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param array|UserModel[] $users
     * @return array
     * @throws Exception
     */
    public function addUsers(array $users): array
    {
        $params = [];
        $values = [];
        $logins = [];
        /** @var UserModel $user */
        foreach ($users as $user) {
            $settings = null;
            if ($user->getSettings() !== null) {
                $settings = json_encode($user->getSettings(), JSON_THROW_ON_ERROR);
            }
            $params[] = $user->getLogin();
            $params[] = $user->getName();
            $params[] = $user->getLastName();
            $params[] = $user->getAge();
            $params[] = $user->getFrom();
            $params[] = $settings;

            $values[] = '(?, ?, ?, ?, ?, ?)';
            $logins[] = $user->getLogin();
        }

        if (count($logins) !== count(array_unique($logins))) {
            throw new RuntimeException('There are duplicates in the user list');
        }

        $exist = count($this->listUsersByLogin($logins));
        if ($exist > 0) {
            throw new RuntimeException($exist . ' users from the list are already registered');
        }

        $connect = $this->connect();
        $connect->beginTransaction();
        try {
            $sth = $connect
                ->prepare($this->getSql('addUsers', implode(',', $values)));
            $sth->execute($params);
            $connect->commit();
        } catch (Exception $e) {
            $connect->rollBack();
            throw $e;
        }

        return $this->listUsersByLogin($logins);
    }
}
