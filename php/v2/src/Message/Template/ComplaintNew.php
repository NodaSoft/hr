<?php

namespace NodaSoft\Message\Template;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationNewInitialData;

class ComplaintNew implements Template
{
    /**
     * @param ReturnOperationNewInitialData $data
     * @param MessageRecipientEntity $recipient
     * @param MessageRecipientEntity $sender
     * @return string
     */
    public function composeSubject(
        InitialData $data,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string {
        return "There is a new complaint ("
            . $data->getMessageTemplate()->getComplaintId()
            . ")";
    }


    /**
     * @param ReturnOperationNewInitialData $data
     * @param MessageRecipientEntity $recipient
     * @param MessageRecipientEntity $sender
     * @return string
     */
    public function composeBody(
        InitialData $data,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string
    {
        $params = $data->getMessageTemplate();
        $message = "Hello, " . $recipient->getFullName()
            . ". Be informed that there is a new complaint №"
            . $params->getComplaintId()
            . "registered. Reseller id: "
            . $data->getReseller()->getId();

        foreach ($params->toArray() as $key => $value) {
            $message .= "$key: $value";
        }

        return $message;
    }
}
