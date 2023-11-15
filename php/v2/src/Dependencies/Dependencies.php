<?php

namespace NodaSoft\Dependencies;

use NodaSoft\Messenger\Client\EmailClient;
use NodaSoft\Messenger\Client\SmsClient;
use NodaSoft\Messenger\Messenger;

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
