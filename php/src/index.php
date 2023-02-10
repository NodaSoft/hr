<?php

require '../vendor/autoload.php';

use Db\PdoConnection;
use Manager\UserManager;

PdoConnection::init("config.php");

$user = new UserManager();


var_dump($user->filterUsersFromAge(10));

$_GET['names'] = ['Тест', 'Вася', '" UNION SELECT(*) FROM users'];
var_dump($user->filterByNames());


$user->addUsers([
    [
        'name' => 'Тест',
        'lastName' => 'Тестов',
        'age' => 18
    ],
    [
        'name' => 'Тест',
        'lastName' => 'Тестов',
        'age' => 'ddd'
    ]
]);

