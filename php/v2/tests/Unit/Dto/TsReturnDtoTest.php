<?php

use NodaSoft\Dto\TsReturnDto;
use PHPUnit\Framework\TestCase;

class TsReturnDtoTest extends TestCase
{
    /** @dataProvider tsReturnDataProvider */
    public function testIsValid(array $data, bool $shouldBeValid): void
    {
        $dto = new TsReturnDto();
        foreach ($data as $key => $value) {
            $setter = 'set' . $key;
            $dto->$setter($value);
        }
        $this->assertSame($shouldBeValid, $dto->isValid());
    }

    public function testToArray(): void
    {
        $data = $this->getValidData();
        $dto = new TsReturnDto();
        foreach ($data as $key => $value) {
            $setter = 'set' . $key;
            $dto->$setter($value);
        }
        $array = $dto->toArray();
        $this->assertSame(count($data), count($array), 'Should be the same number of elements');
        foreach ($data as $key => $value) {
            $keySnake = $this->toUpperSnake($key);
            $this->assertSame($value, $array[$keySnake], 'Should have an analog in SNAKE_CASE');
        }
    }

    public function tsReturnDataProvider(): \Generator
    {
        yield 'valid' => [$this->getValidData(), true];
        yield 'invalid' => [$this->getInvalidData(), false];
    }

    private function getValidData(): array
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
            'differences' => 'Foo Bar Baz',
        ];
    }

    private function getInvalidData(): array
    {
        $data = $this->getValidData();
        array_pop($data); // should be invalid if at least one property absent
        return $data;
    }

    private function toUpperSnake(string $camel): string
    {
        return strtoupper(
            preg_replace('/(?<!^)[A-Z]/', '_$0', $camel)
        );
    }
}
