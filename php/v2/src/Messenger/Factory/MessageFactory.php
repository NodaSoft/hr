<?php

namespace NodaSoft\Messenger\Factory;

use NodaSoft\Messenger\Recipient;
use NodaSoft\Messenger\Message;
use NodaSoft\Messenger\Template\Template;
use NodaSoft\ReferencesOperation\InitialData\InitialData;

class MessageFactory
{
    /** @var Template */
    private $template;

    public function __construct(Template $template)
    {
        $this->template = $template;
    }

    public function makeMessage(
        Recipient   $recipient,
        Recipient   $sender,
        InitialData $initialData
    ): Message {
        $template = $this->template;

        $subject = $template->composeSubject($initialData, $recipient, $sender);
        $body = $template->composeBody($initialData, $recipient, $sender);

        return new Message($content, $recipient, $sender);

        $message->setRecipient($recipient);
        $message->setSender($sender);
        $message->setSubject($subject);
        $message->setBody($body);

        return $message;
    }
}
