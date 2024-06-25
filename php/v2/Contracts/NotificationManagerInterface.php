<?php

namespace NW\WebService\References\Operations\Notification\Contracts;

interface NotificationManagerInterface
{
    public function send(int $resellerId, int $clientId, string $event, int $status, array $templateData, &$error): bool;
}