<?php

namespace App\Enum;

enum Status: int
{
    case COMPLETED = 0;
    case PENDING = 1;
    case REJECTED = 2;

    public static function getNameById(int $id): string
    {
        foreach (self::cases() as $status) {
            if($id === $status->value) {
                return $status->name;
            }
        }

        throw new \ValueError();
    }
}
