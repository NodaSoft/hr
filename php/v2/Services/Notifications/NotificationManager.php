<?php

namespace NW\WebService\References\Operations\Notification\Services\Notifications;

/**
 * Class NotificationManager
 * @package NW\WebService\References\Operations\Notification\Services\Notifications;
 */
class NotificationManager
{
    /**
     * @param int $resellerId
     * @param int $clientId
     * @param string $returnStatus
     * @param int $diffTo
     * @param array $templateData
     * @param string|null $error
     * @return bool
     */
    public static function send(int $resellerId, int $clientId, string $returnStatus, int $diffTo, array $templateData, ?string &$error): bool
    {
        // fake

        return true;
    }
}
