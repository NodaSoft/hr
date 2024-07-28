<?php

namespace NW\WebService\References\Operations\Notification;

/**
 * Contractor class
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;
    public int $id;
    public int $type;
    public string $name;

    /**
     * Constructor for Contractor.
     *
     * @param int $id
     * @param int $type
     * @param string $name
     */
    public function __construct(int $id, int $type = self::TYPE_CUSTOMER, string $name = '')
    {
        $this->id = $id;
        $this->type = $type;
        $this->name = $name;
    }

    /**
     * Get Contractor by ID.
     *
     * @param int $resellerId
     * @return self
     */
    public static function getById(int $resellerId): self
    {
        return new self($resellerId); // fakes the getById method
    }

    /**
     * Get full name of the Contractor.
     *
     * @return string
     */
    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }
}

/**
 * Seller class
 */
class Seller extends Contractor
{
}

/**
 * Employee class
 */
class Employee extends Contractor
{
}

/**
 * Status class
 */
class Status
{
    public int $id;
    public string $name;

    /**
     * Get status name by ID.
     *
     * @param int $id
     * @return string
     */
    public static function getName(int $id): string
    {
        $statuses = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return $statuses[$id] ?? 'Unknown';
    }
}

/**
 * Abstract class for reference operations
 */
abstract class ReferencesOperation
{
    /**
     * Execute the operation.
     *
     * @return array
     */
    abstract public function doOperation(): array;

    /**
     * Get request parameter by name.
     *
     * @param string $pName
     * @return string|null
     */
    public function getRequest(string $pName): ?string
    {
        return $_REQUEST[$pName] ?? null;
    }
}
