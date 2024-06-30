<?php

namespace NW\WebService\References\Operations\Notification\Contracts;

use NW\WebService\References\Operations\Notification\Dto\EmailMessageDto;

interface MessagesClientContract
{
    /**
     * @param EmailMessageDto[] $messages
     * @param int $resellerId
     * @param int $clientId
     * @param int $notificationType
     * @param int|null $differencesTo
     * @return void
     */
    public function sendMessages(
        array $messages,
        int $resellerId,
        int $clientId,
        int $notificationType,
        ?int $differencesTo = null,
    ): void;
}
