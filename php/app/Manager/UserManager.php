<?php

namespace Manager;

if(!defined('IS_INIT'))
    return;

use Core\Db;
use Gateway\UserGateway;

/**
 * Class UserManager
 * @package Manager
 */
class UserManager {
    const limit = 10;

    /**
     * Возвращает пользователей старше заданного возраста.
     * @param int $ageFrom Возраст старше которого выбираем
     * @param int $limit Лимит пользователей
     * @return array
     */
    public static function fetchUsersAgeFrom(int $ageFrom, int $limit=UserManager::limit): array {
        /**
         * Нет смысла делать преобразование типов.
         * Использование trim для чисел, так же сомнительная операция.
         * Название функции не соотвествует её действиям. getUsers больше подойдёт для функции
         * где получают все пользователей
         */
        return UserGateway::getUsersAgeFrom($ageFrom, $limit);
    } // end method

    /**
     * Возвращает пользователей по списку имен.
     * @param array $names Имена пользователей
     * @return array массив пользователей
     */
    public static function getUsersByNames(array $names): array {
        $users = [];
        foreach ($names as $name) {
            if(null !== ($tmp_user = UserGateway::getUserByName($name)))
                $users[] = $tmp_user;
        }
        return $users;
    } // end method

    /**
     * Добавляет пользователей в базу данных.
     * @param $users
     * @return array
     */
    public function saveUsers(array $users): array {
        $ids = [];
        Db::getInstance()->beginTransaction();
        try {
            foreach ($users as $user) {
                if(UserGateway::validate($user))
                    $ids[] = UserGateway::add($user['name'], $user['lastName'], $user['age']);
            }
            Db::getInstance()->commit();
        } catch (\Exception $e) {
            Db::getInstance()->rollBack();
        }
        return $ids;
    } // end method
} // end class