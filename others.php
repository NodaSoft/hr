<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */
class Contractor
{
    //Не указано динамичное свойство класса $resellerId и нет конструктора класса
    const TYPE_CUSTOMER = 0;
    public $id;
    public $type;
    public $name;

    public static function getById(int $resellerId): self
    {
        return new static($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->id ? ($this->name . ' ' . $this->id) : $this->name;
    }
}

class Seller extends Contractor
{
    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
}

class Employee extends Contractor
{
    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
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