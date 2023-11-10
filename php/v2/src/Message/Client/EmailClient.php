<?php

namespace NodaSoft\Message\Client;

use NodaSoft\Message\Client;
use NodaSoft\Message\Message;

class EmailClient implements Client
{
    public function send(Message $message): bool
    {
        return mail(
            $message->getRecipient()->getEmail(),
            $message->getSubject(),
            $message->getBody(),
            $this->getHeaders($message->getSender()->getEmail())
        );
    }

    public function isValid(Message $message): bool
    {
        return $this->isValidEmail($message->getRecipient()->getEmail())
            && $this->isValidEmail($message->getSender()->getEmail())
            && ! empty($message->getBody());
    }

    public function getHeaders(string $emailFrom): string
    {
        return "From: " . $emailFrom;
    }

    private function isValidEmail(?string $email): bool
    {
        return filter_var($email, FILTER_VALIDATE_EMAIL);
    }
}
