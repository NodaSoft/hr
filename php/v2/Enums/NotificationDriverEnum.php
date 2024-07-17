<?php

namespace NW\WebService\References\Operations\Notification\Notification\Enums;

/**
 * Перечисление драйверов уведомлений
 */
enum NotificationDriverEnum: string
{
    case EMAIL = 'email';
    case SMS = 'sms';
}
