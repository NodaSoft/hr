<?php

namespace Nodasoft\Testapp\Notifications;

use Nodasoft\Testapp\Notifications\Base\NotificationInterface;

class SmsNotificationInterfaceClient implements NotificationInterface
{
    public function __construct(array $config = [])
    {
    }

    public function notify(): void
    {
        // some sms logic
    }
}