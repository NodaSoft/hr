<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\DTO\Notification;

/**
 * Class OperationResultDTO
 * Defines Operation Result Data structure
 */
class OperationResultDTO
{
    /**
     * OperationResultDTO Constructor
     *
     * @param bool $isEmployeeNotifiedByEmail
     * @param string|null $employeeEmailNotificationMessage
     * @param bool $isClientNotifiedByEmail
     * @param string|null $clientEmailNotificationMessage
     * @param bool $isClientNotifiedBySms
     * @param string|null $clientSmsNotificationMessage
     */
    public function __construct(
        private bool $isEmployeeNotifiedByEmail = false,
        private ?string $employeeEmailNotificationMessage = null,
        private bool $isClientNotifiedByEmail = false,
        private ?string $clientEmailNotificationMessage = null,
        private bool $isClientNotifiedBySms = false,
        private ?string $clientSmsNotificationMessage = null,
    ) {
    }

    /**
     * @return bool
     */
    public function getIsEmployeeNotifiedByEmail(): bool
    {
        return $this->isEmployeeNotifiedByEmail;
    }

    /**
     * @param bool $value
     * @return void
     */
    public function setIsEmployeeNotifiedByEmail(bool $value): void
    {
        $this->isEmployeeNotifiedByEmail = $value;
    }

    /**
     * @return string|null
     */
    public function getEmployeeEmailNotificationMessage(): ?string
    {
        return $this->employeeEmailNotificationMessage;
    }

    /**
     * @param string $message
     * @return void
     */
    public function setEmployeeEmailNotificationMessage(string $message): void
    {
        $this->employeeEmailNotificationMessage = $message;
    }

    /**
     * @return bool
     */
    public function getIsClientNotifiedByEmail(): bool
    {
        return $this->isClientNotifiedByEmail;
    }

    /**
     * @param bool $value
     * @return void
     */
    public function setIsClientNotifiedByEmail(bool $value): void
    {
        $this->isClientNotifiedByEmail = $value;
    }

    /**
     * @return string|null
     */
    public function getClientEmailNotificationMessage(): ?string
    {
        return $this->clientEmailNotificationMessage;
    }

    /**
     * @param string $message
     * @return void
     */
    public function setClientEmailNotificationMessage(string $message): void
    {
        $this->clientEmailNotificationMessage = $message;
    }

    /**
     * @return bool
     */
    public function getIsClientNotifiedBySms(): bool
    {
        return $this->isClientNotifiedBySms;
    }

    /**
     * @param bool $value
     * @return void
     */
    public function setIsClientNotifiedBySms(bool $value): void
    {
        $this->isClientNotifiedBySms = $value;
    }

    /**
     * @return string|null
     */
    public function getClientSmsNotificationMessage(): ?string
    {
        return $this->clientSmsNotificationMessage;
    }

    /**
     * @param string $message
     * @return void
     */
    public function setClientSmsNotificationMessage(string $message): void
    {
        $this->clientSmsNotificationMessage = $message;
    }
}