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
    public $mobile;
    public $email;

    public static function getById(int $resellerId): self
    {
        return new self($resellerId); // fakes the getById method
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
    const COMPLETED = 0;
    const PENDING = 1;
    const REJECTED = 2;

    public static function getName(int $id): string
    {
        $statusNames = [
            self::COMPLETED => 'Completed',
            self::PENDING => 'Pending',
            self::REJECTED => 'Rejected',
        ];

        return $statusNames[$id] ?? '';
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    protected function getRequest(string $paramName)
    {
        return $_REQUEST[$paramName] ?? null;
    }
}

function getResellerEmailFrom(): string
{
    return 'contractor@example.com';
}

function getEmailsByPermit($resellerId, $event): array
{
    // fakes the method
    return ['someemail@example.com', 'someemail2@example.com'];
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS = 'newReturnStatus';
}
