<?php

namespace NodaSoft\Message\Factory;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\Message\Message;
use NodaSoft\Message\Template\Template;
use NodaSoft\ReferencesOperation\InitialData\TsReturnInitialData;

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
        TsReturnInitialData $initialData
    ): Message {
        $message = new Message();
        $message->setRecipient($recipient);
        $message->setSender($sender);
        $message->setSubject($this->template->composeSubject($initialData));
        $message->setBody($this->template->composeBody($initialData));
        return $message;
    }
}
