<?php

namespace NodaSoft\Message\Client;

use NodaSoft\Message\Message;
use NodaSoft\Message\Client;

class SmsClient implements Client
{
    public function send(Message $message): bool
    {
        //todo: implement sms messenger
    }

    public function isValid(Message $message): bool
    {
        return $this->isCellphoneValid($message->getRecipient()->getCellphone())
            && ! empty($message->getBody());
    }

    public function isCellphoneValid(?int $cellphone): bool
    {
        return strlen((string) $cellphone) === 10;
    }
}
