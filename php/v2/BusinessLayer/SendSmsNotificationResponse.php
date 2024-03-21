<?php

    namespace NW\WebService\References\Operations\Notification\BusinessLayer;

    class SendSmsNotificationResponse
    {


        /** @var bool Была попытка отправки */
        public bool $isSent = false;

        /** @var string Текст ошибка отправки смс */
        public string $message = "";



    }
