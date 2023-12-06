<?php

namespace NW\WebService\References\Operations\Notification\Enums;

abstract class Enum
{
    public static function cases(): array
    {
        return (new \ReflectionClass(static::class))->getConstants();
    }
}