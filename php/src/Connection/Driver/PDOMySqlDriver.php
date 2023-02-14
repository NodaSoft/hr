<?php

namespace App\Connection\Driver;


/**
 * The PDO MySql driver.
 */
class PDOMySqlDriver extends AbstractPDODriver
{
    /**
     * @param array $params
     * @return string
     */
    protected function prepareDns(array $params): string
    {
        return "mysql:dbname={$params['dbname']};host={$params['host']}";
    }

    /**
     * @return string
     */
    public static function getName(): string
    {
        return 'pdo_mysql';
    }
}