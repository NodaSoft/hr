<?php

namespace App\Services;

use App\Enum\Notification;

class MailService
{
    public static function sendMessage(array $data, int $resellerId, string $eventType, ?array $options = null)
    {
        //do something
    }

    public static function getSubject(array $template, int $resellerId): string
    {
        //do somthing
        return '';
    }

    public static function getMessage(array $template, int $resellerId): string
    {
        //do somthing
        return '';
    }
}
