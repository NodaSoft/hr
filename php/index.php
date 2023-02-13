<?php

declare(strict_types=1);

require __DIR__ . '/vendor/autoload.php';

use Envms\FluentPDO;
use App\Domain\User;
use App\Service\UserService;
use App\Infra\UserRepositoryPdo;

// // *** Infra ***
$pdo = new \PDO(
    'mysql:dbname=db;host=127.0.0.1',
    'dbuser',
    'dbpass',
    [
        \PDO::ATTR_ERRMODE => \PDO::ERRMODE_EXCEPTION,
        \PDO::ATTR_DEFAULT_FETCH_MODE => \PDO::FETCH_ASSOC,
    ]
);
$fluent = new FluentPDO\Query($pdo);

$userRepo = new UserRepositoryPdo($fluent);
$userService = new UserService($userRepo);

// *** Usage ***
$userService->getUserByNames(['Ivan', 'Vasya']);
$userService->getUsersFromAge(25);
$userService->addUsers([
    new User(
        firstName: 'Ivan',
        lastName: 'Ivanov',
        age: 28,
        location: 'Moscow',
    ),
    new User(
        firstName: 'Vasya',
        lastName: 'Ivanov',
        age: 24,
        location: 'Kazan',
    ),
]);