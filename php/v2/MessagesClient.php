<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Contracts\MessagesClientInterface;

class MessagesClient implements MessagesClientInterface
{
    public function sendMessage(array $messages, int $resellerId, string $event, int $clientId = null, int $status = null)
    {
        // Реализация отправки сообщений (например, через почту или другой сервис)
    }
}