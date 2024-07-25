<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

enum NotificationEventsEnum: string
{
    case CHANGE_RETURN_STATUS = 'changeReturnStatus';
    case NEW_RETURN_STATUS = 'newReturnStatus';
}