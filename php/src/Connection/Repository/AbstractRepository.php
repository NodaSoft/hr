<?php

namespace App\Connection\Repository;


use App\Connection\DatabaseInterface;

abstract class AbstractRepository implements RepositoryInterface
{
    /**
     * @var DatabaseInterface
     */
    private DatabaseInterface $database;

    /**
     * @param DatabaseInterface $database
     */
    public function __construct(DatabaseInterface $database)
    {
        $this->database = $database;
    }

    /**
     * @return \PDO
     */
    public function getConnection(): \PDO
    {
        return $this->database->getConnection();
    }
}