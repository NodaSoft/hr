<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Contracts;


/**
 * Class Status.
 *
 */
class Status
{
    private const POSIBLE_TYPES = [
        0 => 'Completed',
        1 => 'Pending',
        2 => 'Rejected',
    ];
    /**
     * @var int
     */
    private int $id;
    /**
     * @var string
     */
    private string $name;

    /**
     * get Status name by id
     *
     * @param int $id
     * @return string|null
     */
    public static function getName(int $id): ?string
    {
        return self::POSIBLE_TYPES[$a[$id]] ?? null;
    }
}