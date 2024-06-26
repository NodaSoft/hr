<?php

namespace Nodasoft\Testapp\Enums;

enum NotificationEventType: string
{
    case CHANGE_RETURN_STATUS = 'changeReturnStatus';
    case NEW_RETURN_STATUS = 'newReturnStatus';
}