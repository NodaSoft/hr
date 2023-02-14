<?php

namespace App\Connection\Driver;


/**
 * The PDO MySql driver.
 */
abstract class AbstractPDODriver implements DriverInterface
{
    /**
     * @param array $params
     * @param $username
     * @param $password
     * @return \PDO
     * @throws \InvalidArgumentException
     */
    public function connect(array $params, $username = null, $password = null): \PDO
    {
        try {
            $dns = $this->prepareDns($params);

            return new \PDO($dns, $username, $password);
        } catch (\PDOException $e) {
            throw new \InvalidArgumentException($e->getMessage(), $e->getCode(), $e);
        }
    }

    abstract protected function prepareDns(array $params): string;
}