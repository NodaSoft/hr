<?php

declare(strict_types=1);

namespace NW\WebService\Response;


class ResultOperation
{
    private bool $notificationEmployeeByEmail = false;
    private bool $notificationClientByEmail = false;
    private bool $notificationClientBySms = false;
    private ?string $notificationClientBySmsMessage = null;


    public function setNotificationClientBySmsMess(string $mess): void
    {
        $this->notificationClientBySmsMessage = $mess;
    }

    public function getNotificationClientBySmsMess(): ?string
    {
        return $this->notificationClientBySmsMessage;
    }

    public function notificationEmployeeByEmailSent(): void
    {
        $this->notificationEmployeeByEmail = true;
    }

    public function isNotificationEmployeeByEmailSent(): bool
    {
        return $this->notificationEmployeeByEmail;
    }

    public function notificationClientByEmailSent(): void
    {
        $this->notificationClientByEmail = true;
    }

    public function isNotificationClientByEmailSent(): bool
    {
        return $this->notificationClientByEmail;
    }


    public function notificationClientBySmsSent(): void
    {
        $this->notificationClientBySms = true;
    }

    public function isNotificationClientBySmsSent(): bool
    {
        return $this->notificationClientBySms;
    }


}