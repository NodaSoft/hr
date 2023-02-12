<?php

require_once './vendor/autoload.php';

use \App\Entity\User;

$userManager = new \App\Manager\User();

try {

    $ids = $userManager->addUsers(
        (new User())->setName('User 1')->setAge(21)->setSettings([]),
        (new User())->setName('User 2')->setAge(23),
        (new User())->setName('User 3')->setAge(19),
        (new User())->setName('User 4')->setAge(20)
    );

    var_dump($ids);

} catch(\App\Exception\EntityManagerException $e) {
    exit($e->getMessage());
}

if($names = filter_input(INPUT_GET, 'names', FILTER_DEFAULT, FILTER_REQUIRE_ARRAY)) {
    var_dump($userManager->getByNames($names));
}

var_dump($userManager->getUsers(2, 5));
