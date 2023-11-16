<?php

namespace NodaSoft\Result\Notification;

use NodaSoft\Messenger\Recipient;

class NotificationResult
{
    /** @var bool */
    private $isSent = false;

    /** @var string */
    private $errorMessage = "";

    /** @var Recipient */
    private $recipient;

    public function isSent(): bool
    {
        return $this->isSent;
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

    public function getRecipient(): Recipient
    {
        return $this->recipient;
    }

    public function setRecipient(Recipient $recipient): void
    {
        $this->recipient = $recipient;
    }
}
