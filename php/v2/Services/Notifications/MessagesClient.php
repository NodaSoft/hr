<?php

namespace NW\WebService\References\Operations\Notification\Services\Notifications;

/**
 * Class MessagesClient
 * @package NW\WebService\References\Operations\Notification\Services\Notifications
 */
class MessagesClient
{
    /**
     * @param array $data
     * @param int $resellerId
     * @param string $returnStatus
     * @param int|null $clientId
     * @param int|null $diffTo
     * @return void
     */
    public static function sendMessage(array $data, int $resellerId, string $returnStatus, ?int $clientId = null, ?int $diffTo = null): void
    {
        // fake
    }
}
