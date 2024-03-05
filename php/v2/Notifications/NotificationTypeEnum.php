<?php

namespace NW\WebService\Notification;

enum NotificationTypeEnum: int
{
    case NEW = 1;
    case CHANGE = 2;
}
