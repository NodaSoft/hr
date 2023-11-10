<?php

namespace NodaSoft\Message;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;

class Result
{
    /** @var bool */
    private $isSent;

    /** @var string */
    private $errorMessage;

    /** @var MessageRecipientEntity */
    private $recipient;

    /**
     * @var string
     */
    private $clientClass;

    public function __construct(
        MessageRecipientEntity $recipient,
        string                 $clientClass,
        bool                   $isSent = false,
        string                 $errorMessage = ""
    ) {
        $this->recipient = $recipient;
        $this->clientClass = $clientClass;
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

    public function setRecipient(MessageRecipientEntity $recipient): void
    {
        $this->recipient = $recipient;
    }

    public function getRecipient(): MessageRecipientEntity
    {
        return $this->recipient;
    }

    public function toArray(): array
    {
        return [
            'isSent' => $this->isSent,
            'clientClass' => $this->clientClass,
            'errorMessage' => $this->errorMessage,
            'recipient' => $this->recipient->toArray(),
        ];
    }

    public function getClientClass(): string
    {
        return $this->clientClass;
    }
}
