<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;
    public int $id;
    public string $type;
    public string $name;
    private array $fakes = [
        //...
    ];

    /**
     * @throws \Exception
     */
    public function __construct(int $subjectId)
    {
        $fake = $this->fakes[$subjectId] ?? null;

        if ($fake == null) {
            throw new \Exception();
        }

        $this->id = $fake['id'];
        $this->name = $fake['name'];
    }

    public static function getById(int $resellerId): self
    {
        /**
         * Я бы фейкер сделал. Поскольку в 37 строке будет ошибка. Как минимум рандомом брал бы какие то данные. Или с какого то статического файла.
         * Я только пишу не делаю потому что это займет много времени, и нету смысла.
         */
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

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}