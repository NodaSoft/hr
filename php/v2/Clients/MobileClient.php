<?php

namespace NW\WebService\References\Operations\Notification\Clients;

/**
 * Mobile client for sending notifications.
 *
 * @property NotificationManagerAlias $Messages
 */
class MobileClient implements ClientInterface
{

    public function send(...$args): bool
    {
        return NotificationManagerAlias::send(...$args);
    }
}