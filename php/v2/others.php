<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $seller
 */
class Contractor
{
    public const TYPE_CUSTOMER = 0;
    public int $id;
    public int $type;
    public string $name;

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }

    /**
     * @param int $id
     * @return Contractor
     */
    public function setId(int $id): Contractor
    {
        $this->id = $id;
        return $this;
    }

    /**
     * @return int
     */
    public function getType(): int
    {
        return $this->type;
    }

    /**
     * @param int $type
     * @return Contractor
     */
    public function setType(int $type): Contractor
    {
        $this->type = $type;
        return $this;
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }

    /**
     * @param string $name
     * @return Contractor
     */
    public function setName(string $name): Contractor
    {
        $this->name = $name;
        return $this;
    }

    /**
     * @param int $resellerId
     * @return static
     */
    public static function getById(int $resellerId)
    {
        return new static($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }

    public function isCustomer(): bool
    {
        return $this->getType() === self::TYPE_CUSTOMER;
    }

    public function getSeller(): Seller
    {
        return $this->seller;
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
    public int $id;
    public string $status;

    public static function getStatus(int $id): string
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

function getResellerEmailFrom(): string
{
    return 'contractor@example.com';
}

function getEmailsByPermit($resellerId, $event): array
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}

class NotificationManager
{
    public static function send(
        $resellerId,
        $clientid,
        $event,
        $notificationSubEvent,
        $templateData,
        &$errorText,
        $locale = null
    ) {
        // fakes the method
        return true;
    }
}

class MessagesClient
{
    static function sendMessage(
        $sendMessages,
        $resellerId = 0,
        $customerId = 0,
        $notificationEvent = 0,
        $notificationSubEvent = ''
    ) {
        return '';
    }
}