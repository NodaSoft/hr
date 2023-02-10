<?php
namespace DB;

use PDO;

class DBInstance
{
    /**
     * @var PDO
     */
    public static PDO $instance;

    /**
     * Реализация singleton
     * @return PDO
     */
    public static function getInstance(): PDO
    {
        if (is_null(self::$instance)) {
            $config = require_once 'config/db.php';

            $dsn = sprintf('%s:host=%s;port=%s;dbname=%s',
                $config['database_driver'] ?? 'mysql',
                $config['host'] ?? '127.0.0.1',
                $config['port'] ?? 3306,
                $config['db_name'] ?? 'db'
            );

            $user = $config['username'] ?? 'dbuser';
            $password = $config['password'] ?? 'dbpass';

            self::$instance = new PDO($dsn, $user, $password);
        }

        return self::$instance;
    }

    /**
     * Начать транзакцию
     */
    public static function beginTransaction()
    {
        self::getInstance()->beginTransaction();
    }

    /**
     * commit транзакции
     */
    public static function commit()
    {
        self::getInstance()->commit();
    }

    /**
     * rollback транзакции
     */
    public static function rollback()
    {
        self::getInstance()->rollback();
    }

    /**
     * последний вставленный id
     */
    public static function lastInsertId()
    {
        return self::getInstance()->lastInsertId();
    }
}