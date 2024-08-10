<?php

namespace app\Domain\Notification\Models;

class Status
{
    protected $id, $name;

    public static function getStatus(int $id): string
    {
        $a = [
            100 => 'Completed',
            50 => 'Pending',
            90 => 'Rejected',
        ];

        return $a[$id];
    }
}