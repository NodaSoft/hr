<?php

namespace NW\WebService\References\Operations\Notification\Dto;

readonly class EmailMessageDto
{
    public function __construct(
        public string $emailFrom,
        public string $emailTo,
        public string $subject,
        public string $message,
    ) {
    }
}
