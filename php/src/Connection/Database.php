<?php

namespace App\Connection;


use App\Connection\Driver\DriverInterface;

class Database implements DatabaseInterface
{
    /**
     * @var DriverInterface
     */
    private DriverInterface $driver;

    /**
     * @var array
     */
    private array $params;

    /**
     * @var \PDO
     */
    private \PDO $connection;

    /**
     * @param array $params
     * @param DriverInterface $driver
     */
    public function __construct(\PDO $connection)
    {
        $this->connection = $connection;
    }

    public static function createWithDriver(array $params, DriverInterface $driver): Database
    {
        $username = $params['username'] ?? null;
        $password = $params['password'] ?? null;

        $connection = $driver->connect($params, $username, $password);

        $connection->setAttribute(\PDO::ATTR_ERRMODE, \PDO::ERRMODE_EXCEPTION);
        $connection->setAttribute(\PDO::ATTR_DEFAULT_FETCH_MODE, \PDO::FETCH_ASSOC);

        return new self($connection);
    }

    /**
     * @return \PDO
     */
    public function getConnection(): \PDO
    {
        return $this->connection;
    }
}