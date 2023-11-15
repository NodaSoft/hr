<?php

namespace Tests\Unit\Dto;

use NodaSoft\GenericDto\Dto\ReturnOperationStatusChangedMessageBodyList;
use PHPUnit\Framework\TestCase;

class ReturnOperationStatusChangedMessageBodyListTest extends TestCase
{
    /** @dataProvider tsReturnDataProvider */
    public function testIsValid(array $data, bool $shouldBeValid): void
    {
        $dto = new ReturnOperationStatusChangedMessageBodyList();
        foreach ($data as $key => $value) {
            $setter = 'set' . $key;
            $dto->$setter($value);
        }
        $this->assertSame($shouldBeValid, $dto->isValid());
    }

    public function testToArray(): void
    {
        $data = self::getValidData();
        $dto = new ReturnOperationStatusChangedMessageBodyList();
        foreach ($data as $key => $value) {
            $setter = 'set' . $key;
            $dto->$setter($value);
        }
        $array = $dto->toArray();
        $this->assertSame(count($data), count($array), 'Should be the same number of elements');
        foreach ($data as $key => $value) {
            $keySnake = self::toUpperSnake($key);
            $this->assertSame($value, $array[$keySnake], 'Should have an analog in SNAKE_CASE');
        }
    }

    public function tsReturnDataProvider(): \Generator
    {
        yield 'valid' => [self::getValidData(), true];
        yield 'invalid' => [$this->getInvalidData(), false];
    }

    public static function getValidData(): array
    {
        return [
            'complaintId' => 186,
            'complaintNumber' => 'B234-123',
            'creatorId' => 19894,
            'creatorName' => 'Jake Fillips',
            'expertId' => 17,
            'expertName' => 'Sam Smith',
            'clientId' => 1899,
            'clientName' => 'Tom Sojer',
            'consumptionId' => 50009,
            'consumptionNumber' => 'M654JG',
            'agreementNumber' => 'FF123-4',
            'date' => '2004-11-07',
            'currentStatus' => "reopened",
            'previousStatus' => "closed",
        ];
    }

    private function getInvalidData(): array
    {
        $data = self::getValidData();
        array_pop($data); // should be invalid if at least one property absent
        return $data;
    }

    public static function toUpperSnake(string $camel): string
    {
        return strtoupper(
            preg_replace('/(?<!^)[A-Z]/', '_$0', $camel)
        );
    }
}
