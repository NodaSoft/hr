<?php

namespace NW\WebService\References\Operations\Notification\models;


/**
 * @property ?Seller $Seller
 */
class Contractor
{
    public ?Seller $Seller;
    public string $email;
    const TYPE_CUSTOMER = 0;
    public bool $mobile;

    public function __construct(
        public int $id,
        public int $type = self::TYPE_CUSTOMER,
        public string $name = 'Default Name',
        public ?Seller $seller = null
    ) {
    }
    public static function getById(int $resellerId): ?self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return sprintf('%s %s', $this->name, $this->id);
    }
}