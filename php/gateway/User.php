<?php

namespace Gateway;

use PDO;

class User
{

    /**
     * @var PDO
     */
    private static $instance;

    /**
     * Реализация singleton
     * @return PDO
     */
    public static function getInstance(): PDO
    {
        if (is_null(self::$instance)) {
            // Настройки остаются здесь только для упрощения. В пром. исполнении им здесь не место.
            $dsn = 'mysql:dbname=db;host=127.0.0.1;charset=utf8';
            $user = 'dbuser';
            $password = 'dbpass';
            $options = [
                PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
                PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC, // несмотря на данную настройку, нигде по коду не убираю явное назначение режима, т.к. это не ошибка.
                PDO::ATTR_EMULATE_PREPARES => false,
            ];
            self::$instance = new PDO($dsn, $user, $password, $options);
        }

        return self::$instance;
    }

    /**
     * Возвращает список пользователей старше заданного возраста.
     * @param int $ageFrom
     * @return array
     */
    public static function getUsers(int $ageFrom): array
    {
        $stmt = self::getInstance()->prepare("SELECT id, name, lastName, userFrom, age, settings FROM Users WHERE age > :ageFrom LIMIT " . \Manager\User::LIMIT); //Допускаю формирование запроса в части LIMIT через конкатенацию, т.к. в значении константа класса. В ином случае требует назначения через bindValue c указанием типа
        $stmt->execute(['ageFrom' => $ageFrom]);
        $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);
        $users = [];
        foreach ($rows as $row) {
            $settings = json_decode($row['settings']); //Тип поля в бд?
            $users[] = [
                'id' => $row['id'],
                'name' => $row['name'],
                'lastName' => $row['lastName'],
                'from' => $row['userFrom'],
                'age' => $row['age'],
                'key' => $settings['key'],
            ];
        }

        return $users;
    }

    /**
     * Возвращает пользователя по имени.
     * @param string $name
     * @return array
     */
    public static function user(string $name): array
    {
        $stmt = self::getInstance()->prepare("SELECT id, name, lastName, userFrom, age FROM Users WHERE name = :name");
        $stmt->execute(['name' => $name]);
        $user_by_name = $stmt->fetch(PDO::FETCH_ASSOC);

        return [
            'id' => $user_by_name['id'],
            'name' => $user_by_name['name'],
            'lastName' => $user_by_name['lastName'],
            'from' => $user_by_name['userFrom'], //если заменить ключ с from на userFrom, то можно будет просто возвращать результат от $stmt->fetch(PDO::FETCH_ASSOC);
            'age' => $user_by_name['age'],
        ];
    }

    /**
     * Добавляет пользователя в базу данных.
     * @param string $name
     * @param string $lastName
     * @param int $age
     * @return string
     */
    public static function add(string $name, string $lastName, int $age): string
    {
        $sth = self::getInstance()->prepare("INSERT INTO Users (name, lastName, age) VALUES (:name, :age, :lastName)");
        $sth->execute([':name' => $name, ':age' => $age, ':lastName' => $lastName]);

        return self::getInstance()->lastInsertId();
    }

}

