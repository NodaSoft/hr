<?php

namespace Manager;

use Model\UserModel;

final class UserManager
{

    const MAX_NAMES_FOR_FILTERING = 100;

    /**
     * Возвращает пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    function filterUsersFromAge(int $ageFrom): array
    {
        return UserModel::filterUsersFromAge($ageFrom);
    }

    /**
     * Возвращает пользователей по списку имен.
     * @return array
     */
    public static function filterByNames(): array
    {
        if (!array_key_exists('names', $_GET)) {
            return [];
        }

        $names = $_GET['names'];

        $names = array_slice($names, 0, self::MAX_NAMES_FOR_FILTERING);

        return UserModel::filterByNames($names);
    }

    /**
     * Добавляет пользователей в базу данных.
     * @param $users
     * @return array
     */
    public function addUsers($users): array
    {
        $ids = [];

        try {
            UserModel::getPdoInstance()->beginTransaction();
            foreach ($users as $user) {
                $ids[] = UserModel::add($user['name'], $user['lastName'], $user['age']);
            }
            UserModel::getPdoInstance()->commit();
        } catch (\Exception $e) {
            UserModel::getPdoInstance()->rollBack();

            throw $e;
        }

        return $ids;
    }
}
