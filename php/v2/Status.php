<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

use Exception;

class Status
{
    public const STATUS_COMPLETED = 'Completed';
    public const STATUS_PENDING = 'Pending';
    public const STATUS_REJECTED = 'Rejected';
    public const  STATUSES = [
        0 => self::STATUS_COMPLETED,
        1 => self::STATUS_PENDING,
        2 => self::STATUS_REJECTED,
    ];
    public int $id;

    /**
     * @throws Exception
     */
    public static function getName(int $id): string
    {
        return array_key_exists($id,
            self::STATUSES) ? self::STATUSES[$id] : throw new Exception('Not found status: '.$id);
    }
}
