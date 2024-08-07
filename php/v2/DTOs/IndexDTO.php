<?php

namespace NW\WebService\References\Operations\Notification\DTOs;

class IndexDTO
{
    /**
     * @param bool $notificationEmployeeByEmail
     * @param bool $notificationClientByEmail
     * @param SmsDTO $notificationClientBySms
     */
    public function __construct(
        public bool $notificationEmployeeByEmail,
        public bool $notificationClientByEmail,
        public SmsDTO $notificationClientBySms
    ) {
    }
}
