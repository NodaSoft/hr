<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Interfaces;

/**
 * Interface ContractorInterface
 * Describes Contractor Entity
 */
interface ContractorInterface
{
    /** @var int */
    public const TYPE_CUSTOMER = 1;
    public const TYPE_SELLER = 2;
    public const TYPE_EMPLOYEE = 3;

    /**
     * @return int
     */
    public function getId(): int;

    /**
     * @return string
     */
    public function getFullName(): string;

    /**
     * @return int
     */
    public function getType(): int;

    /**
     * @return string
     */
    public function getEmail(): string;

    /**
     * @return string|null
     */
    public function getMobile(): ?string;
}