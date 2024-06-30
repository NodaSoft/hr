<?php

namespace NW\WebService\References\Operations\Notification;

class Contractor
{
    public const TYPE_CUSTOMER = 0;
    public int $id;
    public int $type;
    public string $name;
    public string $lastName;
    public string $email;
    public string $mobile;
    public Seller $seller;

    public function __construct(int $resellerId)
    {
        $this->id = $resellerId;
    }

    public static function getById(int $resellerId): ?self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return trim($this->name . ' ' . $this->lastName);
    }
}

class Seller extends Contractor
{
    public int $id;
}

class Employee extends Contractor
{
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName)
    {
        return $_REQUEST[$pName];
    }
}