<?php

namespace NW\WebService\References\Operations\Notification;

function getResellerEmailFrom(int $resellerId): string
{
    return 'contractor@example.com';
}

function getEmailsByPermit(int $resellerId, string $event): array
{
    // Faking the method
    return ['someemail@example.com', 'someemail2@example.com'];
}
