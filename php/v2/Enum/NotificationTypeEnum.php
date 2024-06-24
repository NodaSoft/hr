<?php

namespace NW\WebService\References\Operations\Notification;

enum NotificationTypeEnum: int
{
    case TYPE_NEW    = 1;
    case TYPE_CHANGE = 2;
}
