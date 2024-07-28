<?php

namespace NW\WebService\References\Operations\Notification\Notification;

/**
 * Result class
 */
class Result
{
    /**
     * Initializes result array.
     *
     * @return array
     */
    public function initialize(): array
    {
        return [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];
    }
}
