<?php

namespace Nodasoft\Testapp\Notifications;

use Nodasoft\Testapp\Notifications\Base\NotificationInterface;

class MailNotificationInterfaceEmployee implements NotificationInterface
{
    public function __construct(array $config = [])
    {
    }

    public function notify():void
    {

    }
}