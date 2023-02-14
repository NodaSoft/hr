<?php

namespace App\Connection;


use App\Connection\Driver\DriverInterface;
use App\Connection\Driver\PDOMySqlDriver;
use App\Connection\Driver\PDOSqliteDriver;

class DriverFactory
{
    public static function create(array $params): DriverInterface
    {
        $driverName = $params['driver'] ?? '';

        switch ($driverName) {
            case PDOMySqlDriver::getName():
                return new PDOMySqlDriver();
            case PDOSqliteDriver::getName():
                return new PDOSqliteDriver();
        }

        throw new \Exception(sprintf("Driver '%s' not found", $driverName));
    }
}