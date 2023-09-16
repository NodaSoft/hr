<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Struct;

final class Result
{
    /** @phpstan-ignore-next-line */
    private bool $notificationEmployeeByEmail;

    /** @phpstan-ignore-next-line */
    private bool $notificationClientByEmail;

    /** @phpstan-ignore-next-line */
    private SmsNotification $notificationClientBySms;

    public function __construct(
        bool $notificationEmployeeByEmail = false,
        bool $notificationClientByEmail = false,
        SmsNotification $notificationClientBySms = new SmsNotification()
    ) {
        $this->notificationEmployeeByEmail = $notificationEmployeeByEmail;
        $this->notificationClientByEmail   = $notificationClientByEmail;
        $this->notificationClientBySms     = $notificationClientBySms;
    }

    public function notifiedEmployeeByEmail(): self
    {
        $this->notificationEmployeeByEmail = true;
        return $this;
    }

    public function notifiedClientByEmail(): self
    {
        $this->notificationClientByEmail = true;
        return $this;
    }

    public function setSmsNotification(SmsNotification $smsNotification): self
    {
        $this->notificationClientBySms = $smsNotification;
        return $this;
    }
}
