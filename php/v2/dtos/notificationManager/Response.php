<?php

namespace NW\WebService\References\Operations\Notification\notificationManager;

final readonly class Response
{
    public function __construct(
        public ?Error $error = null,
    )
    {
    }
}