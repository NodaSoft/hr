<?php

namespace System\Database;

use PDO;

class PdoProvider
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
        if (!(self::$instance instanceof PDO)) {
            $dsn = "mysql:dbname=db;host=127.0.0.1";
            $user = "dbuser";
            $password = "dbpass";
            self::$instance = new PDO($dsn, $user, $password);
        }

        return self::$instance;
    }
}
