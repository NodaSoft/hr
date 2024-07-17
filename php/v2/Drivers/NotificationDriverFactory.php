<?php

namespace NW\WebService\References\Operations\Notification\Drivers;

use InvalidArgumentException;
use NW\WebService\References\Operations\Notification\Clients\EmailClient;
use NW\WebService\References\Operations\Notification\Clients\MobileClient;
use NW\WebService\References\Operations\Notification\Notification\Enums\NotificationDriverEnum;

/**
 * Class NotificationDriverFactory
 *
 * Фабрика драйверов уведомлений. Выбирает драйвер в зависимости от типа уведомления
 */
class NotificationDriverFactory
{
    public function make(NotificationDriverEnum $driver): NotificationDriverInterface
    {
        return match ($driver) {
            NotificationDriverEnum::EMAIL => new EmailNotificationDriver(new EmailClient()),
            NotificationDriverEnum::SMS => new SmsNotificationDriver(new MobileClient()),
            default => throw new InvalidArgumentException("Driver [$driver->value] not supported."),
        };
    }
}