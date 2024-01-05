<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

class Status
{
    private static $textStatus = [
        'Completed',
        'Pending',
        'Rejected'
    ];

    /**
     * @param int $id
     * @return string
     * @throws Exception
     */
    public static function getName(int $id): string
    {
        if (!key_exists($id, static::$textStatus)) {
            throw new Exception("Status id = {$id} not declared");
        }

        return static::$textStatus[$id];
    }
}
