<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperationResponse
{
    public bool $notificationEmployeeByEmail = false;
    public bool $notificationClientByEmail = false;
    public array $notificationClientBySms = [
        'isSent' => false,
        'message' => '',
    ];

    public function toArray(): array
    {
        return [
            'notificationEmployeeByEmail' => $this->notificationEmployeeByEmail,
            'notificationClientByEmail' => $this->notificationClientByEmail,
            'notificationClientBySms' => [
                'isSent' => $this->notificationClientBySms['isSent'],
                'message' => $this->notificationClientBySms['message'],
            ],
        ];
    }
}
