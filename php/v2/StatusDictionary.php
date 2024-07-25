<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

class StatusDictionary
{
    public static function getName(int $key): string
    {
        $statuses = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return $statuses[$key];
    }
}
