<?php

namespace NW\WebService\References\Operations\Notification;

class MessagesClient
{
    /**
     * Send email to client
     *
     * @param array $messages
     * @param int $resellerId
     * @param string $event
     * @param int|null $clientId
     * @param int|null $to
     * @return void
     */
    public static function sendMessage(array $messages, int $resellerId, string $event, ?int $clientId = null, ?int $to = null): void
    {

    }
}