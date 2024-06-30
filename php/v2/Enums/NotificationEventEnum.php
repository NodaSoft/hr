<?php

namespace NW\WebService\References\Operations\Notification\Enums;

enum NotificationEventEnum: string
{
    case New = 'changeReturnStatus';
    case Change = 'newReturnStatus';
}
