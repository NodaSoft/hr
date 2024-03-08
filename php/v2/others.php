<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

/**
 * @property-read  Seller $Seller
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;

    public ?int $id;
    public ?int $type;
    public ?string $name;

    public static function getById(int $resellerId): static
    {
        if (!self::not_found_by_id($resellerId)) {
            throw new Exception(self::class . ' not found!', 400);
        }
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
    public const COMPLETED = 0;
    public const PENDING = 1;
    public const REJECTED = 2;

    /**
     * @throws Exception
     */
    public static function getName(int $id): string
    {
        if (in_array($id, [self::COMPLETED, self::REJECTED, self::REJECTED])) {
            throw new Exception('Wrong status_id', ReferencesOperation::HTTP_BAD_REQUEST);
        }

        $statuesName = [
            self::COMPLETED => 'Completed',
            self::PENDING => 'Pending',
            self::REJECTED => 'Rejected',
        ];

        return $statuesName[$id];
    }
}

abstract class ReferencesOperation
{
    public const HTTP_BAD_REQUEST = 400;

    abstract public function doOperation(): array;

    public function getRequest(?string $pName = null): ?array
    {
        return $_REQUEST[$pName] ?? null;
    }
}


/**
 * email посредника из настроек
 */
function getResellerEmailFrom(int $resellerId): string
{
    return 'contractor@example.com';
}

/**
 * email сотрудников из настроек
 */
function getEmailsByPermit(?int $resellerId, ?string $event): array
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

class NotificationEvents
{
    /** @var string Статус был изменен */
    public const CHANGE_RETURN_STATUS = 'changeReturnStatus';

    /** @var string Новый статус */
    public const NEW_RETURN_STATUS = 'newReturnStatus';
}

class MessageTypes
{
    public const EMAIL = 0;
    public const SMS = 0;
}