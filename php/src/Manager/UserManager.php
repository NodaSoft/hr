<?php

namespace App\Manager;


use App\Connection\Tools\ConverterTools;
use App\Repository\UserRepository;

class UserManager
{
    const LIMIT = 10;

    private UserRepository $userRepository;

    public function __construct(UserRepository $userRepository)
    {
        $this->userRepository = $userRepository;
    }

    /**
     * Возвращает пользователей старше заданного возраста.
     *
     * @param int $ageFrom
     * @return array
     * @throws \Exception
     */
    function getUsersOlderThanAge(int $ageFrom): array
    {
        $users = $this->userRepository->listByAge($ageFrom, self::LIMIT);

        return array_map([$this, 'prepareUser'], $users);
    }

    /**
     * Возвращает пользователей по списку имен.
     *
     * @param array $names
     * @return array
     * @throws \Exception
     */
    public function listByNames(array $names): array
    {
        $names = array_filter($names, function ($name) {
            return is_string($name);
        });

        $names = array_unique($names);

        $users = [];

        foreach ($names as $name) {
            $users = array_merge($users, $this->userRepository->getByName($name));
        }

        return array_map([$this, 'prepareUser'], $users);
    }

    /**
     * Добавляет пользователей в базу данных.
     *
     * @param array $users
     * @return array
     */
    public function addUsers(array $users): array
    {
        $connection = $this->userRepository->getConnection();

        $connection->beginTransaction();

        $ids = [];

        try {
            foreach ($users as $user) {
                $ids[] = $this->userRepository->add($user['name'], $user['lastName'], $user['age']);
            }

            $connection->commit();
        } catch (\Exception $e) {
            $connection->rollBack();
        }

        return $ids;
    }

    /**
     * @param array $user
     * @return array
     */
    private function prepareUser(array $user): array
    {
        foreach ($user as $k => $v) {
            switch ($k) {
                case 'id':
                case 'age':
                    $v = (int)$v;
                    break;
                case 'settings':

                    $v = ConverterTools::jsonDecode((string)$v);
                    break;
            }

            $user[ $k ] = $v;
        }

        return $user;
    }
}