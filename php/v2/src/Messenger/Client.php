<?php

namespace NodaSoft\Messenger;

interface Client
{
    public function send(
        Message $message,
        Recipient $recipient,
        Recipient $sender
    ): bool;

    public function isValid(
        Message $message,
        Recipient $recipient,
        Recipient $sender
    ): bool;
}
