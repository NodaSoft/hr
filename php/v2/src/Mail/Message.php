<?php

namespace NodaSoft\Mail;

use NodaSoft\DataMapper\EntityInterface\EmailEntity;

class Message
{
    /** @var string */
    private $subject;

    /** @var string */
    private $message;

    /** @var string */
    private $headers;

    /** @var string */
    private $params;

    /** @var EmailEntity */
    private $recipient;

    public function getTo(): string
    {
        return $this->recipient->getEmail();
    }

    public function getSubject(): string
    {
        return $this->subject;
    }

    public function setSubject(string $subject): void
    {
        $this->subject = $subject;
    }

    public function getMessage(): string
    {
        return $this->message;
    }

    public function setMessage(string $message): void
    {
        $this->message = $message;
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

    public function setRecipient(EmailEntity $recipient): void
    {
        $this->recipient = $recipient;
    }

    public function getRecipient(): EmailEntity
    {
        return $this->recipient;
    }
}
