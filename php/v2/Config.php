<?php

namespace NW\WebService\References\Operations\Notification;

class Config
{
    public static function getResellerEmailFrom(int $resellerId): string
    {
        return 'contractor@example.com';
    }

    public static function getEmailsByPermit(int $resellerId, string $event): array
    {
        // fakes the method
        return ['someemeil@example.com', 'someemeil2@example.com'];
    }

}