<?php

namespace Db;

use PDO;

class Db
{
    /**
     * @var static|null
     */
    private static ?Db $instance = null;
    /**
     * @var \PDO|null
     */
    private ?PDO $pdo = null;

    public static function i(): self
    {
        return self::$instance ?? (self::$instance = new static());
    }

    public static function command(string $sql, array $params = []): Command
    {
        return new Command($sql, $params);
    }

    public function pdo(): PDO
    {
        return $this->pdo ?? ($this->pdo = new PDO('mysql:dbname=db;host=127.0.0.1', 'dbuser', 'dbpass'));
    }
}
