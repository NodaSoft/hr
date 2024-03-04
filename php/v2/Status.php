<?php

namespace NW\WebService\References\Operations\Notification;

class Status
{
    //Не используются
    //public $id, $name;

    //Возможно не стоит создавать класс,а вынеси в константы
    public static function getName(int $id): string
    {
        $a = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return $a[$id];
    }
}