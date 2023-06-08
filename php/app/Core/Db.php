<?php
namespace Core;

if(!defined('IS_INIT'))
    return;

class Db extends \PDO {
    private static $instance = null;

    public static function getInstance() {
        if(!self::$instance)
            self::$instance = new Db(DB_DSN, DB_USERNAME, DB_PASSWORD);
        return self::$instance;
    }
}