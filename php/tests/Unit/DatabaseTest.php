<?php

namespace App\Test;

use App\Connection\Database;
use App\Connection\Driver\PDOMySqlDriver;
use App\Connection\DriverFactory;
use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;

class DatabaseTest extends TestCase
{
    /**
     * @var MockObject | Database
     */
    private MockObject $object;

    protected function setUp(): void
    {
        $this->object = $this->createMock(Database::class);
    }

    public function testPdoDriverFactory()
    {
        $result = DriverFactory::create([
            'driver' => PDOMySqlDriver::getName()
        ]);

        $this->assertInstanceOf(PDOMySqlDriver::class, $result);
    }

    public function testGetConnection()
    {
        $result = $this->object->getConnection();

        $this->assertInstanceOf(\PDO::class, $result);
    }
}