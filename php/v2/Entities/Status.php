<?php

namespace NW\WebService\References\Operations\Notification\Entities;

use NW\WebService\References\Operations\Notification\Exceptions\StatusDoesNotExistException;

class Status
{
    public const COMPLETED = 0;
    public const PENDING = 1;
    public const REJECTED = 2;

    private static array $statusNames = [
      self::COMPLETED => 'Completed',
      self::PENDING => 'Pending',
      self::REJECTED => 'Rejected',
    ];

    public static function getName(int $id): string
    {
        return self::$statusNames[$id] ?? throw new StatusDoesNotExistException();
    }
}