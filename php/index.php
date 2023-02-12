<?php

require_once './vendor/autoload.php';

use \App\Entity\User;

$userManager = new \App\Manager\User();

// Добавляение пользователей в базу данных.
try {

    $ids = $userManager->addUsers(
        (new User())->setName('User 1')->setAge(21)->setSettings([]),
        (new User())->setName('User 2')->setAge(23),
        (new User())->setName('User 3')->setAge(19),
        (new User())->setName('User 4')->setAge(20)
    );

    var_dump($ids);

} catch(\Exception $e) {
    exit($e->getMessage());
}

// Получаем пользователей по списку имен.
if($names = filter_input(INPUT_GET, 'names', FILTER_DEFAULT, FILTER_REQUIRE_ARRAY)) {
    var_dump($userManager->getByNames($names));
}

// Получаем пользователей старше заданного возраста.
var_dump($userManager->getUsers(21, 5));
