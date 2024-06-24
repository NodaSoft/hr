<?php

namespace NW\WebService\References\Operations\Notification;

class Contractor
{
    const TYPE_CUSTOMER = 0;
    public $id;
    public $type;
    public $name;

    public static function getById(int $resellerId): ?self
    {
        // Предположительно, этот метод должен получать данные из базы данных или другого источника
        return new self($resellerId); // Заглушка метода getById
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
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
    public $id, $name;

    public static function getName(int $id): string
    {
        $a = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return $a[$id] ?? 'Unknown';
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName)
    {
        return $_REQUEST[$pName] ?? null;
    }
}

function getResellerEmailFrom(int $resellerId): string
{
    return 'contractor@example.com'; // Заглушка метода
}

function getEmailsByPermit(int $resellerId, string $event): array
{
    return ['someemail@example.com', 'someemail2@example.com']; // Заглушка метода
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS = 'newReturnStatus';
}
