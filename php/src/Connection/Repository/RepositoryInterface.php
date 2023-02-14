<?php

namespace App\Connection\Repository;


use App\Connection\DatabaseInterface;

interface RepositoryInterface
{
    public function __construct(DatabaseInterface $database);

    public function getConnection(): \PDO;
}