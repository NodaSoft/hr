<?php

namespace NW\WebService\References\Operations\Notification\Contracts;

use NW\WebService\References\Operations\Notification\Dto\EmailMessageDto;

interface NotificationManagerContract
{
    public function send(
        int $resellerId,
        int $clientId,
        int $notificationType,
        int $differencesTo,
        array $templateData,
    ): bool;
}

