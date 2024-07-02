<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace Models;

use NW\WebService\References\Operations\Notification\Interfaces\ContractorInterface;

/**
 * Class Contractor
 * Model for Entity
 */
class Contractor implements ContractorInterface
{
    /**
     * @param int $id
     * @param string $name
     * @param string $email
     * @param int $type
     * @param string|null $mobile
     */
    public function __construct(
        public int $id,
        public string $name,
        public string $email,
        public int $type = self::TYPE_CUSTOMER,
        public ?string $mobile = null,
    ) {
    }

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }

    /**
     * @return string
     */
    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }

    /**
     * @return string
     */
    public function getEmail(): string
    {
        return $this->email;
    }

    /**
     * @return string|null
     */
    public function getMobile(): ?string
    {
        return $this->mobile;
    }

    /**
     * @return int
     */
    public function getType(): int
    {
        return $this->type;
    }

    /**
     * @return bool
     */
    public function isCustomer(): bool
    {
        return $this->type == self::TYPE_CUSTOMER;
    }

    /**
     * @return bool
     */
    public function isSeller(): bool
    {
        return $this->type == self::TYPE_SELLER;
    }

    /**
     * @return bool
     */
    public function isEmployee(): bool
    {
        return $this->type == self::TYPE_EMPLOYEE;
    }

    /**
     * @return \NW\WebService\References\Operations\Notification\Interfaces\ContractorInterface|null
     */
    public function getSeller(): ?ContractorInterface
    {
        return $this->getSellerId()
            // mock Seller
            ? new self(
                id: $this->getSellerId(),
                name: 'Seller Name',
                email: 'seller@example.com',
                type: self::TYPE_SELLER,
            )
            : null;
    }

    /**
     * @return int|null
     */
    public function getSellerId(): ?int
    {
        return $this->isCustomer() ? 100 : null;
    }

    /**
     * @param string $event
     * @return string[]
     */
    public function getEmailsByPermit(string $event): array
    {
        return ['someemail@example.com', 'someemail2@example.com'];
    }
}