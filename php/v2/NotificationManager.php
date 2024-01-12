<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

class NotificationManager
{
    /**
     * @param  mixed  $resellerId
     * @param  int  $id
     * @param  string  $notificationEvent
     * @param  int  $status
     * @param  array  $templateData
     * @param  string  $error
     * @return bool
     */
    public static function send(
        mixed $resellerId,
        int $id,
        string $notificationEvent,
        int $status,
        array $templateData,
        string $error
    ): bool {
        return true;
    }
}
