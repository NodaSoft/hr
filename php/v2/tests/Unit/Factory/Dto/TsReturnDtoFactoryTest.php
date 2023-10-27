<?php

namespace Tests\Unit\Factory\Dto;

use NodaSoft\Factory\Dto\TsReturnDtoFactory;
use PHPUnit\Framework\TestCase;
use Tests\Unit\Dto\TsReturnDtoTest;

class TsReturnDtoFactoryTest extends TestCase
{
    public function testMakeTsReturnDto(): void
    {
        $data = TsReturnDtoTest::getValidData();
        $factory = new TsReturnDtoFactory();
        $dto = $factory->makeTsReturnDto($data);
        $array = $dto->toArray();
        foreach ($data as $key => $value) {
            $keySnake = TsReturnDtoTest::toUpperSnake($key);
            $this->assertSame($value, $array[$keySnake]);
        }
    }
}
