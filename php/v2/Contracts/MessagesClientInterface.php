<?php

namespace NW\WebService\References\Operations\Notification\Contracts;

interface MessagesClientInterface
{
    public function sendMessage(array $messages, int $resellerId, string $event, int $clientId = null, int $status = null);
}
