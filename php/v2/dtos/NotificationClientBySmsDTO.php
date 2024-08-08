<?php

namespace NW\WebService\References\Operations\Notification;

final class NotificationClientBySmsDTO
{
    public function __construct(
        public bool   $isSent,
        public string $message,
    )
    {
    }

    public function toArray(): array
    {
        return [
            'isSent'  => $this->isSent,
            'message' => $this->message,
        ];
    }
}