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

    public static function getById(int $resellerId): self
    {
        $instance = new self(); // fakes the getById method
        $instance->id = $resellerId;
        $instance->type = self::TYPE_CUSTOMER; // todo: solve the temporary solution
        return $instance;
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }

    public function __get($name)
    {
        //todo: solve the temporary solution (mocking)
        if ($name === "Seller") {
            return Seller::getById(1);
        }
        throw new \InvalidArgumentException("There is no $name property.");
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

        return $a[$id];
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName)
    {
        return $_REQUEST[$pName];
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

function __(): string
{
    /** todo: replace this mocking function */
    return "FOO BAR BAZ Mocking differences";
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}

class MessagesClient
{
    public static function sendMessage(): void
    {
        //todo: implement logic
    }
}

class NotificationManager
{
    public static function send(): bool
    {
        //todo: implement logic
        return true; // fake logic
    }
}
