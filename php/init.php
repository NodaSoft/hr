<?php
if(defined('IS_INIT'))
    return;

define('IS_INIT' , true);
define('ROOT_PATH', dirname(__FILE__));
define('APP_PATH', ROOT_PATH.DIRECTORY_SEPARATOR.'app');
// init autoload
if(!file_exists(dirname(__FILE__).'/config/config.ini'))
    throw new Exception('Конфиг файл не найден');

$config = parse_ini_file(dirname(__FILE__).'/config/config.ini');

define('DB_DSN', $config['DB_DSN']);
define('DB_USERNAME', $config['DB_USERNAME']);
define('DB_PASSWORD', $config['DB_PASSWORD']);

include ROOT_PATH.DIRECTORY_SEPARATOR.'vendor/autoload.php';

spl_autoload_register(function($className) {
    $className = str_replace("\\", DIRECTORY_SEPARATOR, $className);
    $file = APP_PATH.DIRECTORY_SEPARATOR.$className.'.php';
    if (is_readable($file))
        require_once $file;
});
