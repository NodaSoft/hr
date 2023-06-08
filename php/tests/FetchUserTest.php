<?php

include dirname(dirname(__FILE__)).'/init.php';

use Manager\UserManager;

class FetchUserTest extends \PHPUnit\Framework\TestCase {
    public static function setUpBeforeClass(): void {
        // Create users
        $inst = \Core\Db::getInstance();
        $inst->prepare('DELETE FROM users')->execute();
        $prep = $inst->prepare('INSERT INTO users(name, last_name, age) VALUES(:name, :last,:age)');
        $prep->execute([':name' => 'n1',':last' => 'l1', 'age' => 10]);
        $prep->execute([':name' => 'n2',':last' => 'l2', 'age' => 20]);
        $prep->execute([':name' => 'n3',':last' => 'l3', 'age' => 30]);
        parent::setUpBeforeClass();
    }

    public function testFetchByAge() {
        $this->assertEquals(1,count(UserManager::fetchUsersAgeFrom(20)));
        $this->assertEquals(2,count(UserManager::fetchUsersAgeFrom(10)));
        $this->assertEquals(3,count(UserManager::fetchUsersAgeFrom(-1)));
    }

    public function testFetchByName() {
        $users = UserManager::getUsersByNames(['n1','n2','errName']);
        $this->assertEquals(2,count($users));
    }

    public function testInjection() {
        UserManager::getUsersByNames(['"n";DELETE FROM `users`;']);
        $this->assertEquals(3,count(UserManager::fetchUsersAgeFrom(-1)));
    }

}
