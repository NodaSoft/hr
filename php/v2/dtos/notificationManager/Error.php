<?php

namespace NW\WebService\References\Operations\Notification\notificationManager;

final readonly class Error
{
    public function __construct(
        public string $message,
    )
    {
    }
}