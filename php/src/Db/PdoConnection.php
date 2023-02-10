<?php

namespace Db;

use PDO;

class PdoConnection
{

    private static ?PDO $pdoInstance = null;

    private const CONFIG_FILE = "config.php";

    private static function loadConfig($configFile)
    {
        if (!is_readable($configFile)) {
            throw new \Exception(sprintf('Config file %s is not found', $configFile));
        }
        return include_once($configFile);
    }


    public static function init($configFile)
    {
        $config = self::loadConfig($configFile);

        self::$pdoInstance = new PDO($config['dsn'], $config['user'], $config['password']);

        self::$pdoInstance->setAttribute(PDO::ATTR_EMULATE_PREPARES, false);
    }

    /**
     * Реализация singleton
     * @return PDO
     */
    //ну это не синглтон, а хранение экземляра PDO. синглтон хранит экзепляпляр на тот же класс - на Gateway/User в нашем случае
    public static function getPdoInstance(): PDO
    {
        if (is_null(self::$pdoInstance)) {
            throw new \Exception("Use init() at first call");
        }

        return self::$pdoInstance;
    }
}