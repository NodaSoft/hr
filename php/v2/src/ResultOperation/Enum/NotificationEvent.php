<?php

namespace ResultOperation\Enum;

enum NotificationEvent: string
{
    case NEW = 'newReturnStatus';
    case CHANGE = 'changeReturnStatus';
}
