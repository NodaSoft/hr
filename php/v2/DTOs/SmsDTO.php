<?php

namespace NW\WebService\References\Operations\Notification\DTOs;

class SmsDTO
{
    /**
     * @param bool $isSent
     * @param string $message
     */
    public function __construct(
        public bool $isSent,
        public string $message,
    ) {
    }
}
