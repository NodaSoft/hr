<?php

namespace NodaSoft\Dependencies;

use NodaSoft\Message\Client\EmailClient;
use NodaSoft\Message\Client\SmsClient;
use NodaSoft\Message\Messenger;

class Dependencies
{
    public function getMailService(): Messenger
    {
        return new Messenger(new EmailClient());
    }

    public function getSmsService(): Messenger
    {
        return new Messenger(new SmsClient());
    }
}
