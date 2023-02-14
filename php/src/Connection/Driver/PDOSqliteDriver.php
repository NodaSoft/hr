<?php

namespace App\Connection\Driver;


/**
 * The PDO Sqlite driver.
 */
class PDOSqliteDriver extends AbstractPDODriver
{
    /**
     * @param array $params
     * @return string
     */
    protected function prepareDns(array $params): string
    {
        $dsn = 'sqlite:';

        if (isset($params['path'])) {
            $dsn .= $params['path'];
        } elseif (isset($params['memory'])) {
            $dsn .= ':memory:';
        }

        return $dsn;
    }

    /**
     * @return string
     */
    public static function getName(): string
    {
        return 'pdo_sqlite';
    }
}
