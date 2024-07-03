<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;
    public $id;
    public $type;
    public $name;

    public function __construct(int $id, int $type = self::TYPE_CUSTOMER, string $name = '')
    {
        $this->id = $id;
        $this->type = $type;
        $this->name = $name;
    }

    public static function getById(int $id): ?self
    {
        // Здесь должна быть логика получения объекта из базы данных
        // Если объект не найден, возвращаем null
        return new self($id); // Заглушка для метода getById
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
    public $id;
    public $name;

    public static function getName(int $id): string
    {
        $statusNames = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return $statusNames[$id] ?? 'Unknown';
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
    // Здесь должна быть логика получения email по resellerId
    return 'contractor@example.com'; // Заглушка для метода
}

function getEmailsByPermit(int $resellerId, string $event): array
{
    // Здесь должна быть логика получения email'ов по resellerId и событию
    return ['someemail@example.com', 'someemail2@example.com']; // Заглушка для метода
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS = 'newReturnStatus';
}