<?php

namespace NodaSoft\Messenger;

class Result
{
    /** @var bool */
    private $isSent;

    /** @var string */
    private $errorMessage;

    /** @var Recipient */
    private $recipient;

    /** @var string */
    private $clientClass;

    public function __construct(
        Recipient $recipient,
        string $clientClass,
        bool $isSent = false,
        string $errorMessage = ""
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

    public function setRecipient(Recipient $recipient): void
    {
        $this->recipient = $recipient;
    }

    public function getRecipient(): Recipient
    {
        return $this->recipient;
    }

    /**
     * @return array<string, mixed>
     */
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
