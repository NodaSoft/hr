<?php

namespace NW\WebService\References\Operations\Notification;

///**
// * @property Seller $Seller
// */
/**
 * Contractor base class.
 *
 * @property $id
 * @property $type
 * @property $name
 * @property $mobile
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;
    public $id;
    protected $type;
    protected $name;

    protected $mobile;  // добавил

    public static function getById(int $resellerId): self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        // обычно сначала ID, потом NAME...
        return $this->id . ') ' . $this->name;
    }
}

/**
 * Seller extendint the Contractor
 */
class Seller extends Contractor
{
}

/**
 * Employee extendint the Contractor
 */
class Employee extends Contractor
{
}

/**
 * Status class for notification description used as a struct
 *
 * никогда не истанциируется, можно обозначить как abstract,
 * имеет только один static метод.
 */
abstract class Status
{
    // никогда не используются, можно удалить
    // public $id, $name;

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

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}