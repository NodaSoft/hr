<?php

namespace App\Connection;


interface DatabaseInterface
{
    public function getConnection(): \PDO;
}