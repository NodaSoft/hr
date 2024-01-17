<?php


namespace App\v2\Models;

use Exception;

class Status
{
    public const string STATUS_COMPLETED = 'Completed';
    public const string STATUS_PENDING = 'Pending';
    public const string STATUS_REJECTED = 'Rejected';
    public const  array STATUSES = [
        0 => self::STATUS_COMPLETED,
        1 => self::STATUS_PENDING,
        2 => self::STATUS_REJECTED,
    ];
    /**
     * @throws Exception
     */
    public static function getName(int $id): string
    {
        return array_key_exists($id,
            self::STATUSES) ? self::STATUSES[$id] : throw new Exception('Not found status: '.$id);
    }
}
