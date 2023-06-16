<?php

namespace App\Tests;

use App\Service\PdoService;
use App\Service\UserService;
use Exception;

class SimpleTest
{
    private PdoService $pdo;
    private UserService $userService;

    public function __construct()
    {
        $this->pdo = new PdoService();
        $this->userService = new UserService();
    }

    public function run(): void
    {
        // Подготовка
        $this->getPdo()->exe($this->getPdo()->getSql('drop'));
        $this->getPdo()->exe($this->getPdo()->getSql('init'));
        $this->getPdo()->exe($this->getPdo()->getSql('fake'));

        // Тестирование
        $this->getUserAgeFromTest();
        $this->getUserByNameTest();
        $this->getUserByLoginTest();
        $this->listUsersByNameTest();
        $this->listUsersByLoginTest();
        $this->addUserTest();
        $this->addUserExceptionTest();
        $this->addUsersTest();
        $this->addUsersExceptionTest();
    }

    private function getUserAgeFromTest(): void
    {
        $res = $this->getUserService()->getUserAgeFrom(25);
        $right = [
            [
                'id' => 2,
                'login' => 'u2',
                'name' => 'user2',
                'lastName' => 'ln2',
                'from' => 'USA',
                'age' => 27,
                'key' => '222',
            ],
            [
                'id' => 3,
                'login' => 'u3',
                'name' => 'user3',
                'lastName' => 'ln3',
                'from' => null,
                'age' => 35,
                'key' => null,
            ]
        ];
        $this->check('getUserAgeFrom', (string)json_encode($res), (string)json_encode($right));
    }

    private function getUserByNameTest(): void
    {
        $res = $this->getUserService()->getUserByName('user3');
        $right = [
            [
                'id' => 3,
                'login' => 'u3',
                'name' => 'user3',
                'lastName' => 'ln3',
                'from' => null,
                'age' => 35,
            ]
        ];
        $this->check('getUserByName', (string)json_encode($res), (string)json_encode($right));
    }

    private function getUserByLoginTest(): void
    {
        $res = $this->getUserService()->getUserByLogin('u3');
        $right = [
            'id' => 3,
            'login' => 'u3',
            'name' => 'user3',
            'lastName' => 'ln3',
            'from' => null,
            'age' => 35,
        ];
        $this->check('getUserByLogin', (string)json_encode($res), (string)json_encode($right));
    }

    private function listUsersByNameTest(): void
    {
        $res = $this->getUserService()->listUsersByName(['user3', 'user2']);
        $right = [
            [
                'id' => 2,
                'login' => 'u2',
                'name' => 'user2',
                'lastName' => 'ln2',
                'from' => 'USA',
                'age' => 27,
            ],
            [
                'id' => 3,
                'login' => 'u3',
                'name' => 'user3',
                'lastName' => 'ln3',
                'from' => null,
                'age' => 35,
            ]
        ];
        $this->check('listUsersByName', (string)json_encode($res), (string)json_encode($right));
    }

    private function listUsersByLoginTest(): void
    {
        $res = $this->getUserService()->listUsersByLogin(['u3', 'u2']);
        $right = [
            [
                'id' => 2,
                'login' => 'u2',
                'name' => 'user2',
                'lastName' => 'ln2',
                'from' => 'USA',
                'age' => 27,
            ],
            [
                'id' => 3,
                'login' => 'u3',
                'name' => 'user3',
                'lastName' => 'ln3',
                'from' => null,
                'age' => 35,
            ]
        ];
        $this->check('listUsersByLogin', (string)json_encode($res), (string)json_encode($right));
    }

    private function addUserTest(): void
    {
        $res = $this->getUserService()->addUser('tu1', 'test_user1', 'test_ln1', 28);
        $this->check('addUser', (string)$res, '4');
    }

    private function addUserExceptionTest(): void
    {
        try {
            $this->getUserService()->addUser('t', 'test_user1', 'test_ln1', 28);
            $msg = '';
        } catch (Exception $e) {
            $msg = $e->getMessage();
        }
        $this->check('addUser (Exception)', $msg, 'Login must be at least 2 characters.');
    }

    private function addUsersTest(): void
    {
        $right = [
            [
                'id' => 5,
                'login' => 'tu2',
                'name' => 'test_user2',
                'lastName' => 'test_ln2',
                'from' => null,
                'age' => 28
            ],
            [
                'id' => 6,
                'login' => 'tu3',
                'name' => 'test_user3',
                'lastName' => 'test_ln3',
                'from' => null,
                'age' => 38
            ],
        ];

        $res = $this->getUserService()->addUsers(
            [
                ['login' => 'tu2', 'name' => 'test_user2', 'lastName' => 'test_ln2', 'age' => 28],
                ['login' => 'tu3', 'name' => 'test_user3', 'lastName' => 'test_ln3', 'age' => 38],
            ]
        );
        $this->check('addUsers', (string)json_encode($res), (string)json_encode($right));
    }

    private function addUsersExceptionTest(): void
    {
        try {
            $this->getUserService()->addUsers(
                [
                    ['login' => 'tu2', 'name' => 'test_user2', 'lastName' => 'test_ln2', 'age' => 28],
                    ['login' => 'tu2', 'name' => 'test_user3', 'lastName' => 'test_ln3', 'age' => 38],
                ]
            );
            $msg = '';
        } catch (Exception $e) {
            $msg = $e->getMessage();
        }
        $this->check('addUsers (Exception)', $msg, 'There are duplicates in the user list');
    }

    private function check(string $test, string $res, string $right): void
    {
        echo $test, ': ', md5($res) === md5($right) ? 'TRUE' : 'FALSE', PHP_EOL;
    }

    /**
     * @return PdoService
     */
    public function getPdo(): PdoService
    {
        return $this->pdo;
    }

    /**
     * @return UserService
     */
    public function getUserService(): UserService
    {
        return $this->userService;
    }
}
