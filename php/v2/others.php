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

    /**
     * Возвращает экземпляр контрагента по его идентификатору.
     *
     * @param int $resellerId Идентификатор контрагента
     * @return self
     */
    public static function getById(int $resellerId): self
    {
        return new self($resellerId);
    }

    /**
     * Получить полное имя контрагента.
     *
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

class Status
{
    const STATUS_UNKNOWN = 'Неизвестный';

    public $id, $name;

    /**
     * Получить имя статуса по его идентификатору.
     *
     * @param int $id Идентификатор статуса
     * @return string
     */
    public static function getName(int $id): string
    {
        $statusNames = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return $statusNames[$id] ?? self::STATUS_UNKNOWN;
    }
}

abstract class ReferencesOperation
{
    /**
     * Выполнить операцию и вернуть результат.
     *
     * @return array
     */
    abstract public function doOperation(): array;

    /**
     * Получить значение из параметров запроса.
     *
     * @param string $paramName Название параметра
     * @return mixed
     */
    public function getRequest(string $paramName)
    {
        return $_REQUEST[$paramName] ?? null;
    }
}

/**
 * Получить email контрагта, отправителя уведомлений.
 *
 * @param int $resellerId Идентификатор контрагента
 * @return string
 */
function getResellerEmailFrom(int $resellerId): string
{
    return 'contractor@example.com';
}

/**
 * Получить список email-адресов на основе разрешений.
 *
 * @param int $resellerId Идентификатор контрагента
 * @param string $event Событие
 * @return array
 */
function getEmailsByPermit(int $resellerId, string $event): array
{
    return ['someemail@example.com', 'someemail2@example.com'];
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS = 'newReturnStatus';
}
