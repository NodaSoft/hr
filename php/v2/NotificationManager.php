<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Contracts\NotificationManagerInterface;

class NotificationManager implements NotificationManagerInterface
{
    public function send(int $resellerId, int $clientId, string $event, int $status, array $templateData, &$error): bool
    {
        // Реализация отправки уведомлений (например, через SMS или другой сервис)
        return true;
    }
}