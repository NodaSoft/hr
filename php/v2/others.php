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

function getResellerEmailFrom(int $resellerId)
{
    // TODO add get email by reseller id
    return 'contractor@example.com';
}

function getEmailsByPermit(int $resellerId, string $event)
{
    // TODO add get email by reseller id and event
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

function __(string $text,  $templateData, int $resellerId) : string {
    // TODO add normal translate function
    return $text;
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}

class MessagesClient {
    public function sendMessage(array $data, int $resellerId, int $clientId, string $eventType, int $differencesTo = 0 ): void
    {

    }
}

class NotificationManager {
    public function send(int $resellerId, int $clientId, string $eventType, int $differencesTo, $templateData, $error)
    {
        return ['status' => true, 'error' => ''];
    }
}
