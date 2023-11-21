<?php

namespace NW\WebService\References\Operations\Notification;

const TS_GOODS_RETURN = 'tsGoodsReturn';

/**
 * @vphilippov:
 * Не понятно, почему свойство не объявлено в классе.
 * Кроме этого странно, что в базовом классе содержатся объекты дочерних классов.
 *
 * @property Seller $Seller
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;

    public int $id;
    public int $type = self::TYPE_CUSTOMER;
    public string $name = '';
    public ?string $mobile = null;

    public function __construct(int $id)
    {
        $this->id = $id;
    }

    /**
     * @param int $resellerId
     * @return static|null
     */
    public static function getById(int $resellerId): ?static
    {
        return $resellerId ? new static($resellerId) : null; // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
}

class Seller extends Contractor
{
}

/**
 * @vphilippov:
 * Странно, что класс Employee ("сотрудник") наследуется от Contractor ("подрядчик").
 * Логически "сотрудник" не является "подрядчиком". Что-то не так в именовании или логике классов.
 */
class Employee extends Contractor
{
}

class Status
{
    /**
     * @vphilippov:
     * Класс Status везде использутся только статически, нам точно нужны свойства id, name?
     */
    public int $id;
    public string $name;

    const COMPLETED = 'Completed';
    const PENDING   = 'Pending';
    const REJECTED  = 'Rejected';

    public static function getName(int $id): ?string
    {
        return static::all()[$id] ?? null;
    }

    public static function checkId(int $id): bool
    {
        return !is_null(static::all()[$id] ?? null);
    }

    public static function all(): array
    {
        return [
            self::COMPLETED,
            self::PENDING,
            self::REJECTED
        ];
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName): mixed
    {
        return $_REQUEST[$pName] ?? null;
    }
}

function getResellerEmailFrom(): string
{
    return 'contractor@example.com';
}

function getEmailsByPermit(int $resellerId, string $event): array
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}