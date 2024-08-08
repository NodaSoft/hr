<?php

namespace NW\WebService\References\Operations\Notification\components;

use NW\WebService\References\Operations\Notification\NotificationEvents;

class MessagesClient
{
    public function sendMessage(array $data, int $resellerId, int $id, NotificationEvents $event, int $diffTo = 0): void
    {
    }
}