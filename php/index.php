<?php

namespace Anoshin\Managers;

require 'vendor/autoload.php';

$manager = new Manager();
var_dump($manager->getOlderThan(10));

if (!empty($_GET['names'])) {
	var_dump($manager->getByNames($_GET['names']));
}

$managers = [
	[
		'name' => 'Вася',
		'age' => 20,
		'lastName' => 'Иванов'
	],
	[
		'name' => 'Петя',
		'age' => 30,
		'lastName' => 'Петров'
	],
];

var_dump($manager->insertAll($managers));

