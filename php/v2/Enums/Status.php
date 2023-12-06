<?php

namespace NW\WebService\References\Operations\Notification\Enums;

class Status extends Enum
{
    public const COMPLETED = 0;
    public const PENDING = 1;
    public const REJECTED = 2;

    public static function getName(int $id): string
    {
        $statuses = array_flip(self::cases());

        if (!isset($statuses[$id])) {
            throw new \Exception('Status not found!');
        }

        return ucfirst(strtolower($statuses[$id]));
    }
}