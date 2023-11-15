<?php

namespace NodaSoft\Messenger;

use NodaSoft\GenericDto\Dto\Dto;

class Message
{
    /** @var string */
    private $subject;

    /** @var string */
    private $body;

    public function __construct(Content $content, Dto $params)
    {
        $this->subject = $content->composeMessageSubject($params);
        $this->body = $content->composeMessageBody($params);
    }

    public function getSubject(): string
    {
        return $this->subject;
    }

    public function getBody(): string
    {
        return $this->body;
    }
}
