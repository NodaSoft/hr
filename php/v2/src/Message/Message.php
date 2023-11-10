<?php

namespace NodaSoft\Message;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;

class Message
{
    /** @var string */
    private $subject;

    /** @var string */
    private $body;

    /** @var string */
    private $headers;

    /** @var string */
    private $params;

    /** @var MessageRecipientEntity */
    private $recipient;

    /** @var MessageRecipientEntity */
    private $sender;

    public function getSubject(): string
    {
        return $this->subject;
    }

    public function setSubject(string $subject): void
    {
        $this->subject = $subject;
    }

    public function getBody(): string
    {
        return $this->body;
    }

    public function setBody(string $body): void
    {
        $this->body = $body;
    }

    public function getHeaders(): string
    {
        return $this->headers ?? "";
    }

    public function setHeaders(string $headers): void
    {
        $this->headers = $headers;
    }

    public function getParams(): string
    {
        return $this->params ?? "";
    }

    public function setParams(string $params): void
    {
        $this->params = $params;
    }

    public function setRecipient(MessageRecipientEntity $recipient): void
    {
        $this->recipient = $recipient;
    }

    public function getRecipient(): MessageRecipientEntity
    {
        return $this->recipient;
    }

    public function getSender(): MessageRecipientEntity
    {
        return $this->sender;
    }

    public function setSender(MessageRecipientEntity $sender): void
    {
        $this->sender = $sender;
    }
}
