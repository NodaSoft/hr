<?php

declare(strict_types=1);

namespace ResultOperation\Event;

use pseudovendor\Event;
use ResultOperation\DTO\NotificationTemplate;
use ResultOperation\Entity\Contractor;

abstract class AbstractStatusEvent extends Event
{
    /**
     * @var NotificationTemplate
     */
    protected readonly NotificationTemplate $template;

    /**
     * @var Contractor
     */
    protected readonly Contractor $client;

    /**
     * @var string
     */
    protected string $error;

    /**
     * @var bool
     */
    private bool $employeeByEmail = false;

    /**
     * @var bool
     */
    private bool $clientBySms = false;

    /**
     * @var bool
     */
    private bool $clientByEmail = false;

    /**
     * @return NotificationTemplate
     */
    public function getTemplate(): NotificationTemplate
    {
        return $this->template;
    }

    /**
     * @return Contractor
     */
    public function getClient(): Contractor
    {
        return $this->client;
    }

    /**
     * @param string $message
     * @return $this
     */
    public function setError(string $message): self
    {
        $this->error = $message;

        return $this;
    }

    /**
     * @return ?string
     */
    public function getError(): ?string
    {
        return $this->error ?? null;
    }

    /**
     * @return void
     */
    public function setEmployeeByEmail(): void
    {
        $this->employeeByEmail = true;
    }

    /**
     * @return void
     */
    public function setClientBySms(): void
    {
        $this->clientBySms = true;
    }

    /**
     * @return void
     */
    public function setClientByEmail(): void
    {
        $this->clientByEmail = true;
    }

    /**
     * @return bool
     */
    public function isEmployeeNotifiedByEmail(): bool
    {
        return $this->employeeByEmail;
    }

    /**
     * @return bool
     */
    public function isClientNotifiedBySms(): bool
    {
        return $this->clientBySms;
    }

    /**
     * @return bool
     */
    public function isClientNotifiedByEmail(): bool
    {
        return $this->clientByEmail;
    }
}
