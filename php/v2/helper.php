<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

const RESELLER_EMAIL_FROM = 'contractor@example.com';

/**
 * @return string
 */
function getResellerEmailFrom(): string
{
    return RESELLER_EMAIL_FROM;
}

/**
 * @param $resellerId
 * @param $event
 * @return string[]
 */
function getEmailsByPermit($resellerId, $event): array
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}
