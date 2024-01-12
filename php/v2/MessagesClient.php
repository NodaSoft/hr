<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

class MessagesClient
{
    /**
     * @param  array  $parameters
     * @param  mixed  $resellerId
     * @param  string  $notificationEvent
     * @param  int|null  $clientId
     * @param  int|null  $status
     * @return void
     */
    public static function sendMessage(
        array $parameters,
        mixed $resellerId,
        string $notificationEvent,
        int $clientId = null,
        int $status = null
    ): void {
        //TODO::implementation
    }
}
