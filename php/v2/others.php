<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $seller
 */
class Contractor
{
    // Не совсем понятно зачем тут TYPE_CUSTOMER
    const TYPE_CUSTOMER = 0;

    public function __construct(
        public readonly int $id = 0,
        public int $type = static::TYPE_CUSTOMER,
        public string $name = '',
        public string $email = '',
        public string $mobile = '',
    )
    {
    }

    public static function getById(int $resellerId): self
    {
        return new static($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id; // пойдет, можно и лучше
    }

    // Для получения свойства seller
    public function __get($name): mixed
    {
        if ($name === 'seller') {
            return "something";
        } else {
            return null;
        }
    }
}

class Seller extends Contractor
{
}

class Employee extends Contractor
{
}

class Status
{
    private const NAMES = [
        0 => 'Completed',
        1 => 'Pending',
        2 => 'Rejected',
    ];

    public static function getName(int $id): string
    {
        return self::NAMES[$id] ?? 'undefined';
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(array $data): array;

    // этот геттер тут не нужен, кмк
    // надо передавать данные непосредственно в метод как я сделал выше
    public function getRequest($pName)
    {
        return $_REQUEST[$pName];
    }
}

function getResellerEmailFrom(int $resellerId): string
{
    // some logic here
    return 'contractor@example.com';
}

function getEmailsByPermit($resellerId, $event)
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

enum NotificationEvents: string
{
    case CHANGE_RETURN_STATUS = 'changeReturnStatus';
    case NEW_RETURN_STATUS = 'newReturnStatus';
}