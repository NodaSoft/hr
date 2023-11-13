<?php

namespace NodaSoft\Message\Factory;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\Message\Message;
use NodaSoft\Message\Template\Template;
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
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender,
        InitialData $initialData
    ): Message {
        $template = $this->template;

        $subject = $template->composeSubject($initialData, $recipient, $sender);
        $body = $template->composeBody($initialData, $recipient, $sender);

        $message = new Message();
        $message->setRecipient($recipient);
        $message->setSender($sender);
        $message->setSubject($subject);
        $message->setBody($body);

        return $message;
    }
}
