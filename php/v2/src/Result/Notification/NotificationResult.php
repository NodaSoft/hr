<?php

namespace NodaSoft\Result\Notification;

class NotificationResult
{
    /** @var bool */
    private $isSent;

    /** @var string */
    private $errorMessage;

    public function isSent(): bool
    {
        return $this->isSent ?? false;
    }

    public function setIsSent(bool $isSent): void
    {
        $this->isSent = $isSent;
    }

    public function getErrorMessage(): string
    {
        return $this->errorMessage;
    }

    public function setErrorMessage(string $errorMessage): void
    {
        $this->errorMessage = $errorMessage;
    }

    public function toArray(): array
    {
        $array = [];
        foreach ($this as $key => $value) {
            $array[$key] = $value;
        }
        return $array;
    }
}
