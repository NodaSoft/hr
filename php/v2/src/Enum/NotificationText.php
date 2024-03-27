<?php

namespace App\Enum;

enum NotificationText: string
{
    case NEW = 'NewPositionAdded';
    case CHANGED = 'PositionStatusHasChanged';
}
