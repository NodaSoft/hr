<?php

include dirname(dirname(__FILE__)).'/init.php';

use Gateway\UserGateway;

class CreateUserTest extends \PHPUnit\Framework\TestCase {

    public function testCreate() {
        $last_id = UserGateway::add('name', 'lastName', 2);
        $this->assertGreaterThan(0, $last_id);
    }

    public function testValidate() {
        $this->assertEquals(true, UserGateway::validate([
            'name' => 'Name',
            'lastName' => 'Last Name',
            'age' => 10
        ]));
        $this->assertEquals(false, UserGateway::validate([
            'name' => 'Name',
            'age' => 10
        ]));
        $this->assertEquals(false, UserGateway::validate([
            'name' => 'Name',
            'lastName' => 'Last Name',
            'age' => -10
        ]));
    }
}
