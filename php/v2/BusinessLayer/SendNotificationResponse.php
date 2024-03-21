<?php

    namespace NW\WebService\References\Operations\Notification\BusinessLayer;

    class SendNotificationResponse
    {
        public bool $notificationEmployeeByEmail = false;
        public bool $notificationClientByEmail = false;
        public SendSmsNotificationResponse $notificationClientBySms;

        public function __construct() {
            $this->notificationClientBySms = new SendSmsNotificationResponse();
        }

        public function toArray(): array {
            return json_decode(json_encode($this), true);
        }

    }
