<?php

namespace NW\WebService\References\Operations\Notification;

final class ResponseDTO
{
    public function __construct(
        public bool $notificationEmployeeByEmail,
        public bool $notificationClientByEmail,
        public NotificationClientBySmsDTO $notificationClientBySms
    )
    {
    }

    public function toArray(): array
    {
        return [
            'notificationEmployeeByEmail' => $this->notificationEmployeeByEmail,
            'notificationClientByEmail'   => $this->notificationClientByEmail,
            'notificationClientBySms'     => $this->notificationClientBySms->toArray(),
        ];
    }
}