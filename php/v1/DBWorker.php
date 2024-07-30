<?php

declare(strict_types=1);


namespace DB;

use PDO;

class DBWorker
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
            $dsn = 'mysql:dbname=db;host=127.0.0.1';
            $user = 'dbuser';
            $password = 'dbpass';
            self::$instance = new PDO($dsn, $user, $password);
        }

        return self::$instance;
    }

}