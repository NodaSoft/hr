<?php

namespace NW\WebService\References\Operations\Notification\Notification\Enum;

/**
 * Notification event enum
 */
enum NotificationEvent: string
{
    case CHANGE_RETURN_STATUS = 'changeReturnStatus';
    case NEW_RETURN_STATUS = 'newReturnStatus';
}
