<?php

namespace Model;

use PDO;

final class UserModel extends Model
{
    const LIMIT = 10;
    const DEFAULT_SETTINGS = ['key' => null];


    //исходим из того, что тип у поля JSON c NULL
    private static function getUserSettings(?string $userSettings): array
    {
        $settings = [];

        if (is_string($userSettings)) {
            $tmpSettings = json_decode($userSettings, JSON_OBJECT_AS_ARRAY);

            //если в поле лежит не JSON массив, а, например, строка.
            if (is_array($tmpSettings)) {
                $settings = $tmpSettings;
            }
        }

        //склеиваем со стандартными настройками
        return array_merge(self::DEFAULT_SETTINGS, $settings);
    }

    /**
     * Возвращает список пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    public static function filterUsersFromAge(int $ageFrom): array
    {
        //лучше использовать подстановку аргументов, а не склейку строк
        //столбцы должны быть в обратных кавычках, особенно `from`
        $stmt = self::getPdoInstance()->prepare(
            "SELECT `id`, `name`, `lastName`, `from`, `age`, `settings` FROM `Users` WHERE `age` > :age LIMIT :limit"
        );
        $stmt->execute([
            ':age' => $ageFrom,
            ':limit' => self::LIMIT
        ]);
        $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);
        $users = [];

        foreach ($rows as $row) {
            $settings = self::getUserSettings($row['settings']);

            $users[] = [
                'id' => $row['id'],
                'name' => $row['name'],
                'lastName' => $row['lastName'],
                'from' => $row['from'],
                'age' => $row['age'],
                'key' => $settings['key'],
            ];
        }

        return $users;
    }


    private static function prepareList(string|array $list): string
    {
        return implode(
            ',',
            array_map(
                fn(string $a) => self::getPdoInstance()->quote($a),
                (array)$list
            )
        );
    }


    /**
     * Возвращает пользователя по имени.
     * @return array
     */
    public static function filterByName(string $name): array
    {
        $stmt = self::getPdoInstance()->prepare(
            "SELECT `id`, `name`, `lastName`, `from`, `age` FROM `Users` WHERE `name` = :name"
        );
        $stmt->execute([":name" => $name]);

        return $stmt->fetch(PDO::FETCH_ASSOC);
    }

    /**
     * Возвращает пользователей по имени.
     * @param string|list<string> $names
     * @return array
     */
    public static function filterByNames(string|array $names): array
    {
        $names = self::prepareList($names);

        $stmt = self::getPdoInstance()->prepare(
            "SELECT `id`, `name`, `lastName`, `from`, `age` FROM `Users` WHERE `name` IN ($names)"
        );
        $stmt->execute();

        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    /**
     * Добавляет пользователя в базу данных.
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @return false|int
     */
    public static function add(string $name, string $lastName, int $age): false|int
    {
        $sth = self::getPdoInstance()->prepare(
            "INSERT INTO `Users` (`name`, `lastName`, `age`) VALUES (:name, :lastName, :age)"
        );
        //нужна обработка ошибок
        $sth->execute([':name' => $name, ':age' => $age, ':lastName' => $lastName]);

        //может возвращать false, тогда будет несовпадение типов
        //если все в порядке, то стоит возвращать сразу int
        $id = self::getPdoInstance()->lastInsertId();

        if ($id !== false) {
            $id = (int)$id;
        }
        return $id;
    }
}