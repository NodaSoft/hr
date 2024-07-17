<?php

namespace NW\WebService\References\Operations\Notification\Clients;

/**
 * Email client for sending notifications.
 *
 * @property MessagesClient $Messages
 */
class EmailClient implements ClientInterface
{
    public function send(...$args): bool
    {
        return new MessagesClient::sendMessage(...$args);
    }
}