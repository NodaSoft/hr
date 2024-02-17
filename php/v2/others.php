<?php

namespace NW\WebService\References\Operations\Notification;

enum MessageTypes: int
{
    case EMAIL = 0;
}

/**
 * @property Seller $seller
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;
    public $id;
    public $type;
    public $name;

    public static function getById(int $resellerId): self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->id ? ($this->name . ' ' . $this->id) : $this->name;
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
    const COMPLETED = 0;
    const PENDING = 1;
    const REJECTED = 2;
    public $id, $name;

    public static function getName(int $id): string
    {
        $statusValues = [
            self::COMPLETED => 'Completed',
            self::PENDING => 'Pending',
            self::REJECTED => 'Rejected',
        ];

        return $statusValues[$id];
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($paramName)
    {
        return $_REQUEST[$paramName];
    }
}

function getResellerEmailFrom()
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