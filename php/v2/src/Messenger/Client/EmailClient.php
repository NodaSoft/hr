<?php

namespace NodaSoft\Messenger\Client;

use NodaSoft\Messenger\Recipient;
use NodaSoft\Messenger\Client;
use NodaSoft\Messenger\Message;

class EmailClient implements Client
{
    public function send(
        Message $message,
        Recipient $recipient,
        Recipient $sender
    ): bool {
        return mail(
            $recipient->getEmail(),
            $message->getSubject(),
            $message->getBody(),
            $this->getHeaders($sender->getEmail())
        );
    }

    public function isValid(
        Message $message,
        Recipient $recipient,
        Recipient $sender
    ): bool {
        return $this->isValidEmail($recipient->getEmail())
            && $this->isValidEmail($sender->getEmail())
            && ! empty($message->getBody());
    }

    public function getHeaders(string $emailFrom): string
    {
        return "From: " . $emailFrom;
    }

    private function isValidEmail(?string $email): bool
    {
        return (bool) filter_var($email, FILTER_VALIDATE_EMAIL);
    }
}
