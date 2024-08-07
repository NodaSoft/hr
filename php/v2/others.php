<?php

namespace NW\WebService\References\Operations\Notification;

abstract class Contractor
{
    const TYPE_CLIENT = 1;
    const TYPE_SELLER = 2;
    const TYPE_EMPLOYEE = 3;

    private int $id;
    private string $name;
    private string $email;
    private string $mobile;

    public static function getById(int $resellerId): ?self
    {
        return new static($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->getName() . ' ' . $this->getId();
    }

    public function getId(): int
    {
        return $this->id;
    }

    public abstract function getType(): int;

    public function getName(): string
    {
        return $this->name;
    }

    public function getEmail(): string
    {
        return $this->email;
    }

    public function getMobile(): string
    {
        return $this->mobile;
    }
}

class Client extends Contractor
{
    private Seller $seller;
    public function getType(): int
    {
        return static::TYPE_CLIENT;
    }

    public function getSeller(): Seller
    {
        return $this->seller;
    }

}

class Seller extends Contractor
{
    public function getType(): int
    {
        return static::TYPE_SELLER;
    }
}

class Employee extends Contractor
{
    public function getType(): int
    {
        return static::TYPE_EMPLOYEE;
    }
}

class Status
{

    public static function getName(int $id): string
    {
        $statuses = [
          0 => 'Completed',
          1 => 'Pending',
          2 => 'Rejected',
        ];

        return $statuses[$id] ?? $statuses[2];
    }
}

abstract class ReferencesOperation
{

    abstract public function doOperation(): array;

    public function getRequest(string $paramName): array
    {
        return $_REQUEST[$paramName] ?? [];
    }

}

function getResellerEmailFrom(): string
{
    return 'contractor@example.com';
}

function getEmailsByPermit(int $resellerId, string $event): array
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

Enum NotificationEvents: int
{
    case CHANGE_RETURN_STATUS = 1000;
    case NEW_RETURN_STATUS    = 1001;
}
