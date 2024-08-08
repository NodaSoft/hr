<?php

namespace NW\WebService\References\Operations\Notification\components;

use NW\WebService\References\Operations\Notification\NotificationEvents;
use NW\WebService\References\Operations\Notification\notificationManager\Response;

class NotificationManager
{
    public function send(
        int $resellerId,
        int $id,
        NotificationEvents $event,
        int $param,
        array $templateData
    ): Response
    {
        return new Response();
    }
}