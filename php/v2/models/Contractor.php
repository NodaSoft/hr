<?php

namespace NW\WebService\References\Operations\Notification\models;

use NW\WebService\References\Operations\Notification\ContractorType;
use NW\WebService\References\Operations\Notification\Seller;

/**
 * @property Seller $seller
 */
class Contractor
{
    public ContractorType $type;
    public string $name;
    public string $email;
    public string $mobile;

    private function __construct(
        public readonly int $id
    )
    {
    }

    public static function findById(int $id): ?static
    {
        return new static($id); // fakes the findById method
    }

    public function getFullName(): string
    {
        return "$this->name $this->id";
    }
}