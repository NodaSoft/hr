<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */
class Contractor
{
    private int $id;
    public $type;
    private string $name;
    private string $email;
    private string $mobile;

    public static function getById(int $id): self
    {
        return new self($id); // fakes the getById method
    }

    public function getId(): int
    {
        return $this->id;
    }

    public function getFullName(): string
    {
        return trim($this->name . ' ' . $this->id);
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

class Customer extends Contractor
{
}

class Seller extends Contractor
{
}

class Employee extends Contractor
{
}

class Status
{
    // если статусы, не предполагается хранить в БД, то можно было сделать так
    public const STATUS_COMPLETED = 0;
    public const STATUS_PENDING = 1;
    public const STATUS_REJECTED = 2;
    public const STATUSES = [
        self::STATUS_COMPLETED => 'Completed',
        self::STATUS_PENDING => 'Pending',
        self::STATUS_REJECTED => 'Rejected',
    ];

    // но вот из-за этого, я сделал предположение, что статусы хранятся в БД и сделал аналог получения объекта
    public $id, $name;

    public static function getById(int $id): self
    {
        return new self($id); // fakes the getById method
    }

    public function getNameSelf():string
    {
        return $this->name;
    }

    public static function getName(int $id): ?string
    {
        return self::getById($id)?->getNameSelf();
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName): array
    {
        return !empty($_REQUEST[$pName]) && is_array($_REQUEST[$pName]) ? $_REQUEST[$pName] : [];
    }
}

function getSellerEmailFrom()
{
    return 'contractor@example.com';
}

function getEmailsByPermit($resellerId, $event)
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}