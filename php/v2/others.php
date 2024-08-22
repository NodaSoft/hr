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
        $name = $this->name . ' ' . $this->id;
        return strlen(trim($name)) > 0 ? $name : $this->name;
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
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    private static function validateRequest(array $request): array
    {
        $keys = [
            'reseller_id',
            'notificationType',
            'clientId',
            'creatorId',
            'expertId',
            'complaintId',
            'complaintNumber',
            'consumptionId',
            'consumptionNumber',
            'agreementNumber'
        ];
        foreach ($keys as $key) {
            //просто проверим, что они есть, чтобы получить требуемые модели и гарантировано сделать $templateData
            if (!isset($request[$key]) || intval($request[$key]) <= 0 || is_array($request[$key])) {
                throw new \InvalidArgumentException("Incorrect {$key}", 422);
            }
        }
        if ($request['notificationType'] === self::TYPE_CHANGE) {
            foreach (['from', 'to'] as $key) {
                if (!isset($request['differences'][$key]) ||
                    intval($request['differences'][$key]) <= 0 ||
                    is_array($request['differences'][$key])) {
                    throw new \InvalidArgumentException('Incorrect difference id', 422);
                }
            }
        }
        return $request;
    }
    /**
     * @param mixed $value
     * @return mixed
     */
    private static function escapeRequestValue($value)
    {
        if (is_array($value)) {
            foreach ($value as $k => $v) {
                $value[$k] = self::escapeRequestValue($v);
            }
        } elseif (is_string($value)) {
            $value = htmlspecialchars($value, ENT_QUOTES, 'UTF-8');
        }

        return $value;
    }

    abstract public function doOperation(): array;

    protected function getRequest($pName): array
    {
        $request = [];
        if (isset($_REQUEST[$pName])) {
            foreach ($_REQUEST[$pName] as $k => $v) {
                $request[$k] = self::escapeRequestValue($v);
            }
        } else {
            throw new \Exception("Request not contain '{$pName} or empty request'");
        }

        return self::validateRequest($request);
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