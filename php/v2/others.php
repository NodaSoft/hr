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
        /**
         * По идее, метод всегда возвращает объект, поэтому исключения не может быть.
         * Т.к. это "заглушка" и в будущем все же нужны будут исключения, тогда можно либо общее исключения вида:
         * throw new \Exception('Contractor not found!', 400);
         * либо, под каждый потом свой текст прописать. Для этого переназначить этот метод с нужным текстом.
         */
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }

    public function isIdEqualTo(int $id):bool
    {
        return $this->id === $id;

    }
    public function isCustomer():bool
    {
        return $this->type===self::TYPE_CUSTOMER;
    }
}

class Seller extends Contractor
{
}

class Employee extends Contractor
{
}
class Expert extends Contractor
{
}
class Client extends Contractor
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

        return $a[$id]??'';
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    /**
     * @throws \Exception
     */
    public function getRequest($pName):array
    {
        if ($_REQUEST[$pName]) {
            return $_REQUEST[$pName];
        }
        throw new \Exception('Empty request', 400);
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