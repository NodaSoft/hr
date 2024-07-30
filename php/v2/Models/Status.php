<?php

namespace NW\WebService\References\Operations\Notification\Models;

/**
 * Class Status
 * @package NW\WebService\References\Operations\Notification\Models;
 */
class Status
{
    /** @var int $id */
    public int $id;

    /** @var string $name */
    public string $name;

    /**
     * @param int $id
     * @return string
     */
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
