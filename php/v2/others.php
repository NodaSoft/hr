<?php

namespace NW\WebService\References\Operations\Notification;


abstract class Status
{
    const Completed = 0;
    const Pending = 1;
    const Rejected = 2;

    public static function getStatus(int $statusId): string
    {
        $statuses = [
            self::Completed => 'Completed',
            self::Pending => 'Pending',
            self::Rejected => 'Rejected',
        ];

        return $statuses[$statusId] ?? '';
    }
}

function getEmailsByPermit($resellerId, $event)
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}

class NotificationEvents
{
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}

