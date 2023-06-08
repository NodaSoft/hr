<?php

namespace Gateway;

if(!defined('IS_INIT'))
    return;

use Core\Db;
/**
 * Называть классы всё же лучше по разному
 * Class UserGateway
 * @package Gateway
 */
class UserGateway
{
     /**
     * Возвращает список пользователей старше заданного возраста.
     * @param int $ageFrom
     * @param int $limit
     * @return array
     */
    public static function getUsersAgeFrom(int $ageFrom, int $limit): array {
        $stmt = Db::getInstance()->prepare("SELECT * FROM users WHERE age > :age LIMIT :limit");
        $stmt->bindParam(':age',$ageFrom, Db::PARAM_INT);
        $stmt->bindParam(':limit',$limit, Db::PARAM_INT);
        $stmt->execute();
        $rows = $stmt->fetchAll(Db::FETCH_ASSOC);
        $users = [];
        foreach ($rows as $row) {
            /**
             * Дальше с settings работают как с массивом, необходим параметр assoc=true
             * json_decode($row['settings'], true)
             */
            $settings = json_decode($row['settings'], true);
            $users[] = [
                'id' => $row['id'],
                'name' => $row['name'],
                'lastName' => $row['last_name'],
                'from' => $row['from'],
                'age' => $row['age'],
                'key' => $settings['key']??null, // Необходимо учитывать если парамеры нет
            ];
        }
        return $users;
    } // end method

    /**
     * Возвращает пользователя по имени.
     * @param string $name
     * @return mixed array|null
     */
    public static function getUserByName(string $name): ?array {
        $inst = Db::getInstance();
        $stmt = $inst->prepare("SELECT `id`, `name`, `last_name`, `from`, `age`, `settings` FROM users WHERE `name` = :name");
        $stmt->execute([
            'name' => $name
        ]);
        $user_by_name = $stmt->fetch(Db::FETCH_ASSOC);

        // Проверка если пользователя не нашли
        if(!$user_by_name)
            return null;

        $settings = json_decode($user_by_name['settings'], true);
        return [
            'id' => $user_by_name['id'],
            'name' => $user_by_name['name'],
            'lastName' => $user_by_name['last_name'],
            'from' => $user_by_name['from'],
            'age' => $user_by_name['age'],
            'key' => $settings['key']??null, // Необходимо учитывать если парамеры нет
        ];
    } // end method

    /**
     * Добавляет пользователя в базу данных.
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @return int
     */
    public static function add(string $name, string $lastName, int $age): int {
        $sth = Db::getInstance()->prepare("INSERT INTO users (name, last_name, age) VALUES (:name, :lastName, :age)");
        $sth->execute(['name' => $name, 'age' => $age, 'lastName' => $lastName]);
        return Db::getInstance()->lastInsertId();
    } // end method

    /**
     * Простенькая валидация
     * @param array $user
     * @return bool
     */
    public static function validate(array $user) : bool {
        $valid = true;
        $valid &= array_key_exists('age', $user) && $user['age'] > 0;
        $valid &= array_key_exists('name', $user) && preg_match('/\w+/', $user['name']);
        $valid &= array_key_exists('lastName', $user) && preg_match('/\w+/', $user['lastName']);
        return $valid;
    } // end method
} // end class