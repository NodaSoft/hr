<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

class MessagesClient
{
    public static function sendMessage(
        array $params,
        int $resellerId,
        int $clientId,
        ?string $notificationType = null,
        ?int $differencesTo = null
    ): void {

    }
}
