<?php

namespace NodaSoft\Messenger\Client;

use NodaSoft\Messenger\Recipient;
use NodaSoft\Messenger\Message;
use NodaSoft\Messenger\Client;

class SmsClient implements Client
{
    public function send(
        Message $message,
        Recipient $recipient,
        Recipient $sender
    ): bool {
        //todo: implement sms messenger
    }

    public function isValid(
        Message $message,
        Recipient $recipient,
        Recipient $sender
    ): bool {
        return $this->isCellphoneValid($recipient->getCellphone())
            && ! empty($message->getBody());
    }

    public function isCellphoneValid(?int $cellphone): bool
    {
        return strlen((string) $cellphone) === 10;
    }
}
