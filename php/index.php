<?php

require_once './vendor/autoload.php';

use \App\Entity\User;

$userManager = new \App\Manager\User();

/*$ids = $userManager->addUsers(
    (new User())->setName('User 1')->setAge(21)->setSettings([]),
    (new User())->setName('User 2')->setAge(23),
    (new User())->setName('User 3')->setAge(19),
    (new User())->setName('User 4')->setAge(20),
);

var_dump($ids);*/

// var_dump($userManager->getByNames());

var_dump($userManager->getUsers(2));
