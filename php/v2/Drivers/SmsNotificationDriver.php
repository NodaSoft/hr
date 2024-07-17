<?php

namespace NW\WebService\References\Operations\Notification\Drivers;

use NW\WebService\References\Operations\Notification\Clients\ClientInterface;

/**
 * Class SmsNotificationDriver
 *
 *  Обработчик уведомлений по смс.
 */
class SmsNotificationDriver implements NotificationDriverInterface
{
    public function __construct(private ClientInterface $client)
    {
    }

    /**
     * Sends an SMS notification to the recipient.
     */
    public function send(array $data): bool
    {
        return $this->client->send($data);
    }
}