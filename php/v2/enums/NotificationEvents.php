<?php

namespace NW\WebService\References\Operations\Notification;

enum NotificationEvents: string
{
    case CHANGE_RETURN_STATUS = 'changeReturnStatus';
    case NEW_RETURN_STATUS    = 'newReturnStatus';
}