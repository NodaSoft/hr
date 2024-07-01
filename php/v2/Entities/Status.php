<?php

namespace NW\WebService\References\Operations\Notification\Entities;

class Status
{
  public int $id;
  public string $name;

  private static array $statusNames = [
    0 => 'Completed',
    1 => 'Pending',
    2 => 'Rejected',
  ];

  public static function getName(int $id): string
  {
    return self::$statusNames[$id] ?? 'Unknown';
  }
}