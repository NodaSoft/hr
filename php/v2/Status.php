<?php

namespace NW\WebService\References\Operations\Notification;

class Status
{
    public $id;
    public $name;
    
    private static $names = [
        0 => 'Completed',
        1 => 'Pending',
        2 => 'Rejected',
    ];
    
    public static function getNameById( int $id): string
    {
        return self::$names[$id];
    }
}
