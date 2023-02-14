<?php

require '../vendor/autoload.php';

$config = include "../config.php";

$driver = \App\Connection\DriverFactory::create($config);

$database = \App\Connection\Database::createWithDriver($config, $driver);

$userRepository = new \App\Repository\UserRepository($database);

$userManager = new \App\Manager\UserManager($userRepository);

echo "<pre>";

/*
$users = $userManager->addUsers([
    [
        'name'     => 'qwe',
        'lastName' => 'qwe',
        'age'      => 10,
    ],
    [
        'name'     => 'asd',
        'lastName' => 'asd',
        'age'      => 20,
    ],
    [
        'name'     => 'zxc',
        'lastName' => 'zxc',
        'age'      => 30,
    ],
]);
print_r($users);
*/

$users = $userManager->listByNames([
    'qwe',
    'asd',
    'test',
]);
print_r($users);

$users = $userManager->getUsersOlderThanAge(5);
print_r($users);