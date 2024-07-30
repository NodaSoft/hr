<?php

namespace NW\WebService\References\Operations\Notification\Models;

/**
 * Class Contractor
 * @package NW\WebService\References\Operations\Notification\Models
 *
 * @property Seller $Seller
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;

    /** @var int $id */
    public int $id;

    /** @var int $type */
    public int $type;

    /** @var string $name */
    public string $name;

    /**
     * Find by ID or return null
     *
     * @param int $resellerId
     * @return self|null
     */
    public static function getById(int $resellerId): ?self
    {
        return new self($resellerId); // fakes the getById method
    }

    /**
     * @return string
     */
    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
}
