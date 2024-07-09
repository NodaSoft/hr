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
    public $email;
    public $mobile;
    public $Seller;

    public function __construct(int $id = 0)
    {
        $this->id     = $id;
        $this->type   = self::TYPE_CUSTOMER;
        $this->name   = 'ContractorName';
        $this->email  = 'client@example.com';
        $this->mobile = '1234567890';
        $this->Seller = new Seller($id);
    }

    public static function getById(int $id): ?self
    {
        // Эмулируем получение объекта по id
        if ($id > 0) {
            return new self($id);
        }
        return null;
    }

    public function getFullName(): string
    {
        return htmlspecialchars($this->name . ' ' . $this->id, ENT_QUOTES, 'UTF-8');
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
    public $id;
    public $name;

    public static function getName(int $id): string
    {
        $a = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return htmlspecialchars($a[$id], ENT_QUOTES, 'UTF-8');
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName)
    {
        // Добавляем фильтрацию входных данных
        return htmlspecialchars($pName, ENT_QUOTES, 'UTF-8');
    }
}

function getResellerEmailFrom($resellerId)
{
    // Предполагается, что здесь возвращается проверенный email
    return 'contractor@example.com';
}

function getEmailsByPermit($resellerId, $event)
{
    // Эмулируем получение email по разрешению
    return ['someemail@example.com', 'someemail2@example.com'];
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}

class MessagesClient
{
    public static function sendMessage(array $messages, int $resellerId, $clientIdOrEvent, $event = null, $statusTo = null)
    {
        // Эмулируем отправку сообщения
        return true;
    }
}

class NotificationManager
{
    public static function send($resellerId, $clientId, $event, $statusTo, $templateData, &$error)
    {
        // Эмулируем отправку уведомления
        $error = '';
        return true;
    }
}

function __($key, $data = null, $resellerId = null)
{
    // Эмулируем перевод строки
    // Экранируем вывод, чтобы защитить от XSS
    return htmlspecialchars($key, ENT_QUOTES, 'UTF-8');
}
