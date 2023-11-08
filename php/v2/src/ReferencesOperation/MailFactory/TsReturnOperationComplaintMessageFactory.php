<?php

namespace NodaSoft\ReferencesOperation\MailFactory;

use NodaSoft\DataMapper\EntityInterface\EmailEntity;
use NodaSoft\Mail\Message;
use NodaSoft\ReferencesOperation\InitialData\TsReturnInitialData;

class TsReturnOperationComplaintMessageFactory
{
    public function makeMessage(
        EmailEntity $recipient,
        TsReturnInitialData $initialData
    ): Message {
        $resellerEmail = $initialData->getReseller()->getEmail();
        $message = new Message();
        $message->setRecipient($recipient);
        $message->setSubject($this->composeSubject($initialData));
        $message->setMessage($this->composeMessage($initialData));
        $message->setHeaders($this->composeHeaders($resellerEmail));
        return $message;
    }

    public function composeSubject(TsReturnInitialData $initialData): string
    {
        //todo: handle template logic

        return "Complaint claim ("
            . $initialData->getMessageTemplate()->getDate()
            . ")";
    }

    public function composeMessage(TsReturnInitialData $initialData): string
    {
        //todo: handle template logic
        $message = "There is a complaint claim. Reseller id: "
            . $initialData->getReseller()->getId();

        $template = $initialData->getMessageTemplate()->toArray();
        foreach ($template as $key => $value) {
            $message .= "$key: $value";
        }

        return $message;
    }

    public function composeHeaders(string $senderEmail): string
    {
        return "From: $senderEmail";
    }
}
