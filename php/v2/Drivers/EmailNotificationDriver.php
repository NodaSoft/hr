<?php

namespace NW\WebService\References\Operations\Notification\Drivers;

use NW\WebService\References\Operations\Notification\Clients\ClientInterface;

/**
 * Class EmailNotificationDriver
 *
 * Обработчик уведомлений по электронной почте.
 */
readonly class EmailNotificationDriver implements NotificationDriverInterface
{
    public function __construct(private ClientInterface $client)
    {
    }

    /**
     * Sends an email notification to the recipient.
     */
    public function send(array $data): bool
    {
        return $this->client->send($data);
    }
}