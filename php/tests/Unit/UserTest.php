<?php

namespace App\Test;

use App\Connection\Database;
use App\Repository\UserRepository;
use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;

class UserTest extends TestCase
{
    /**
     * @var UserRepository
     */
    private UserRepository $object;

    /**
     * @var MockObject | \PDOStatement
     */
    private MockObject $pdoStatement;

    /**
     * @var MockObject | \PDO
     */
    private MockObject $pdo;

    protected function setUp(): void
    {
        $this->pdo = $this->createMock(\PDO::class);
        $this->pdoStatement = $this->createMock(\PDOStatement::class);

        $database = $this->createMock(Database::class);

        $this->object = new UserRepository($database);

        $database->expects($this->any())
            ->method('getConnection')
            ->willReturn($this->pdo);
    }

    public function testListByAgeEmpty()
    {
        $age = 99;
        $limit = 0;
        $expectedResult = [];

        $this->pdo->expects($this->once())
            ->method('prepare')
            ->with($this->equalTo('SELECT id, name, lastName, `from`, age, settings FROM `Users` WHERE age > :age LIMIT :limit'))
            ->willReturn($this->pdoStatement);

        $this->pdoStatement->expects($this->exactly(2))
            ->method('bindParam')
            ->withConsecutive([
                    $this->equalTo(':age'),
                    $this->equalTo($age),
                ], [
                    $this->equalTo(':limit'),
                    $this->equalTo($limit),
                ],);

        $this->pdoStatement->expects($this->once())
            ->method('execute');

        $this->pdoStatement->expects($this->once())
            ->method('fetchAll')
            ->willReturn($expectedResult);

        $result = $this->object->listByAge($age, $limit);

        $this->assertEmpty($result);
    }

    public function testListByAgeNotEmpty()
    {
        $age = 99;
        $limit = 1;
        $expectedResult = [
            [],
        ];

        $this->pdo->expects($this->once())
            ->method('prepare')
            ->with($this->equalTo('SELECT id, name, lastName, `from`, age, settings FROM `Users` WHERE age > :age LIMIT :limit'))
            ->willReturn($this->pdoStatement);

        $this->pdoStatement->expects($this->exactly(2))
            ->method('bindParam')
            ->withConsecutive([
                    $this->equalTo(':age'),
                    $this->equalTo($age),
                ], [
                    $this->equalTo(':limit'),
                    $this->equalTo($limit),
                ]);

        $this->pdoStatement->expects($this->once())
            ->method('execute');

        $this->pdoStatement->expects($this->once())
            ->method('fetchAll')
            ->willReturn($expectedResult);

        $result = $this->object->listByAge($age, $limit);

        $this->assertNotEmpty($result);
    }

    public function testGetByNameEmpty()
    {
        $name = '';
        $expectedResult = [];

        $this->pdo->expects($this->once())
            ->method('prepare')
            ->with($this->equalTo('SELECT id, name, lastName, `from`, age FROM `Users` WHERE name = :name'))
            ->willReturn($this->pdoStatement);

        $this->pdoStatement->expects($this->exactly(1))
            ->method('bindParam')
            ->with(
                $this->equalTo(':name'), $this->equalTo($name),
            );

        $this->pdoStatement->expects($this->once())
            ->method('execute');

        $this->pdoStatement->expects($this->once())
            ->method('fetchAll')
            ->willReturn($expectedResult);

        $result = $this->object->getByName($name);

        $this->assertEmpty($result);
    }

    public function testGetByNameNotEmpty()
    {
        $name = 'test';
        $expectedResult = [
            [],
        ];

        $this->pdo->expects($this->once())
            ->method('prepare')
            ->with($this->equalTo('SELECT id, name, lastName, `from`, age FROM `Users` WHERE name = :name'))
            ->willReturn($this->pdoStatement);

        $this->pdoStatement->expects($this->exactly(1))
            ->method('bindParam')
            ->with(
                $this->equalTo(':name'), $this->equalTo($name),
            );

        $this->pdoStatement->expects($this->once())
            ->method('execute');

        $this->pdoStatement->expects($this->once())
            ->method('fetchAll')
            ->willReturn($expectedResult);

        $result = $this->object->getByName($name);

        $this->assertNotEmpty($result);
    }
}