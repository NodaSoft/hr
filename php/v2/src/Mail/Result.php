<?php

namespace NodaSoft\Mail;

use NodaSoft\DataMapper\EntityInterface\EmailEntity;

class Result
{
    /** @var bool */
    private $isSent;

    /** @var string */
    private $errorMessage;

    /** @var EmailEntity */
    private $recipient;

    public function __construct(
        EmailEntity $recipient,
        bool $isSent = false,
        string $errorMessage = ""
    ) {
        $this->recipient = $recipient;
        $this->isSent = $isSent;
        $this->errorMessage = $errorMessage;
    }

    public function setIsSent(bool $isSent): void
    {
        $this->isSent = $isSent;
    }

    public function isSent(): bool
    {
        return $this->isSent;
    }

    public function setErrorMessage(string $errorMessage): void
    {
        $this->errorMessage = $errorMessage;
    }

    public function getErrorMessage(): string
    {
        return $this->errorMessage;
    }

    public function setRecipient(EmailEntity $recipient): void
    {
        $this->recipient = $recipient;
    }

    public function getRecipient(): EmailEntity
    {
        return $this->recipient;
    }

    public function toArray(): array
    {
        return [
            'isSent' => $this->isSent,
            'errorMessage' => $this->errorMessage,
            'recipient' => $this->recipient->toArray(),
        ];
    }
}
