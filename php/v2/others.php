<?php
// PHP 8.1
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
    private bool $mobile = false;
    public ?string $email = null;
    /**
     * @param int $resellerId
     *
     * @return static|null
     */
    public static function getById(int $resellerId): ?self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }

    public function hasFullName(): bool
    {
        return !empty($this->name);
    }

    public function isMobile(): bool
    {
        return $this->mobile;
    }

    public function hasEmail():bool
    {
        return !empty($this->email);
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
    private static $statusEnum = [
        0 => 'Completed',
        1 => 'Pending',
        2 => 'Rejected',
    ];
    /**
     * @param int $id
     *
     * @return string
     * @throws \Exception
     */
    public static function getName(int $id): string
    {
        if(!isset(self::$statusEnum[$id]))
        {
            throw new \Exception("Неверный статус");
        }

        return self::$statusEnum[$id];
    }

    public static function isValid(int $id): bool
    {
        return isset(self::$statusEnum[$id]);
    }
}

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName)
    {
        return $_REQUEST[$pName] ?? null;
    }
}

function getResellerEmailFrom(int $resellerId): ?string
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