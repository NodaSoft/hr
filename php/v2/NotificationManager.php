<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

class NotificationManager
{
    public static function send(int $resellerId, int $clientId, string $notification, int $differenceTo, array $templateData, &$error): array
    {
        return [];
    }
}
