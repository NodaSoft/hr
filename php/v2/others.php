<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 * @property string $email
 * @property string mobile
 */
class Contractor
{
    public const int TYPE_CUSTOMER = 0;
    private int $id;

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }

    /**
     * @param int $id
     * @return void
     */
    public function setId(int $id): void
    {
        $this->id = $id;
    }

    /**
     * @var int
     */
    private int $type;

    /**
     * @return int
     */
    public function getType(): int
    {
        return $this->type;
    }

    /**
     * @param int $type
     * @return void
     */
    public function setType(int $type): void
    {
        $this->type = $type;
    }

    /**
     * @var string
     */
    private string $name;

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }

    /**
     * @param string $name
     * @return void
     */
    public function setName(string $name): void
    {
        $this->name = $name;
    }

    /**
     * @param int $resellerId
     * @return self
     */
    public static function getById(int $resellerId): self
    {
        return new self($resellerId); // fakes the getById method
    }

    /**
     * @return string
     */
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

enum Status: int
{
    case Completed = 0;
    case Pending = 1;
    case Rejected = 3;
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName): array
    {
        return isset($_REQUEST[$pName]) ? (array)($_REQUEST[$pName]) : [];
    }
}

function getResellerEmailFrom(): string
{
    return 'contractor@example.com';
}

function getEmailsByPermit($resellerId, $event): array
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

enum NotificationEvents: string
{
    case CHANGE_RETURN_STATUS = 'changeReturnStatus';
    case NEW_RETURN_STATUS = 'newReturnStatus';
}