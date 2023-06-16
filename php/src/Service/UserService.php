<?php

declare(strict_types=1);

namespace App\Service;

use App\Model\UserModel;
use App\Repository\UserRepository;
use Exception;
use RuntimeException;

class UserService
{
    private const LIMIT = 10;

    private readonly UserRepository $userRepo;

    public function __construct()
    {
        $this->userRepo = new UserRepository();
    }

    /**
     * Возвращает пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    public function getUserAgeFrom(int $ageFrom): array
    {
        return $this->getUserRepo()->getUserAgeFrom($ageFrom, self::LIMIT);
    }

    /**
     * @param string $names
     * @return array|null
     */
    public function getUserByName(string $names): ?array
    {
        return $this->getUserRepo()->getUserByName($names, self::LIMIT);
    }

    public function getUserByLogin(string $login): ?array
    {
        return $this->getUserRepo()->getUserByLogin($login);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @param array $names
     * @return array
     */
    public function listUsersByName(array $names): array
    {
        return $this->getUserRepo()->listUsersByName($names);
    }

    public function listUsersByLogin(array $logins): array
    {
        return $this->getUserRepo()->listUsersByLogin($logins);
    }

    private function checkUser(UserModel $user): void
    {
        if (mb_strlen($user->getLogin()) < 2) {
            throw new RuntimeException('Login must be at least 2 characters.');
        }
        if (mb_strlen($user->getName()) < 1) {
            throw new RuntimeException('Name must be at least 1 characters.');
        }

        if (mb_strlen($user->getLastName()) < 1) {
            throw new RuntimeException('Last name must be at least 1 characters.');
        }

        if ($user->getAge() < 18 || $user->getAge() > 150) {
            throw new RuntimeException('Age must be over 18 and under 150.');
        }
        // ...
    }

    /**
     * Добавляет пользователя в БД
     * @param string $login
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @param string|null $from
     * @param array|null $settings
     * @return int
     */
    public function addUser(
        string $login,
        string $name,
        string $lastName,
        int $age,
        string $from = null,
        array $settings = null
    ): int {
        $user = (new UserModel())
            ->setLogin($login)
            ->setName($name)
            ->setLastName($lastName)
            ->setAge($age)
            ->setFrom($from)
            ->setSettings($settings);
        $this->checkUser($user);
        return $this->getUserRepo()->addUser($user);
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param array $rawUsers
     * @return array
     * @throws Exception
     */
    public function addUsers(array $rawUsers): array
    {
        $users = [];
        foreach ($rawUsers as $rawUser) {
            $user = $this->userMapping($rawUser);
            try {
                $this->checkUser($user);
            } catch (Exception $e) {
                throw new Exception($e->getMessage() . ' ' . json_encode($rawUser, JSON_THROW_ON_ERROR));
            }
            $users[] = $user;
        }
        return $this->getUserRepo()->addUsers($users);
    }

    private function userMapping(array $rawUser): UserModel
    {
        $settings = null;
        if (isset($rawUser['settings']) && is_array($rawUser['settings'])) {
            $settings = $rawUser['settings'];
        }
        return (new UserModel())
            ->setLogin(isset($rawUser['login']) ? (string)$rawUser['login'] : '')
            ->setName(isset($rawUser['name']) ? (string)$rawUser['name'] : '')
            ->setLastName(isset($rawUser['lastName']) ? (string)$rawUser['lastName'] : '')
            ->setAge(isset($rawUser['age']) ? (int)$rawUser['age'] : 0)
            ->setFrom((isset($rawUser['from']) && is_string($rawUser['from']) ? $rawUser['from'] : null))
            ->setSettings($settings);
    }

    /**
     * @return UserRepository
     */
    public function getUserRepo(): UserRepository
    {
        return $this->userRepo;
    }
}
