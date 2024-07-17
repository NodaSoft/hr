<?php

namespace NW\WebService\References\Operations\Notification\Notification\Enums;

/**
 * Перечисление типов уведомлений
 */
enum NotificationTypeEnum: int
{
    case NEW = 1;
    case CHANGE = 2;
}
