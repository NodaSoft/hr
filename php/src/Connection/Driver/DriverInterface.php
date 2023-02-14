<?php

namespace App\Connection\Driver;

/**
 * DriverInterface
 */
interface DriverInterface
{
    public function connect(array $params, $username = null, $password = null): \PDO;

    public static function getName(): string;
}