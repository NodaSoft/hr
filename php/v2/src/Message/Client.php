<?php

namespace NodaSoft\Message;

interface Client
{
    public function send(Message $message): bool;

    public function isValid(Message $message): bool;
}
