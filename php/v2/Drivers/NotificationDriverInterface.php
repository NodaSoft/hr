<?php

namespace NW\WebService\References\Operations\Notification\Drivers;

/**
 * Interface NotificationDriverInterface
 *
 * Defines the contract for notification drivers.
 */
interface NotificationDriverInterface
{
    /**
     * Sends a notification to the recipient.
     */
    public function send(array $data): bool;
}
