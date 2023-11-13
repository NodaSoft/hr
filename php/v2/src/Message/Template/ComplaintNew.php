<?php

namespace NodaSoft\Message\Template;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationNewInitialData;

class ComplaintNew implements Template
{
    /**
     * @param ReturnOperationNewInitialData $initialData
     * @param MessageRecipientEntity $recipient
     * @param MessageRecipientEntity $sender
     * @return string
     */
    public function composeSubject(
        InitialData $initialData,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string {
        return "There is a new complaint ("
            . $initialData->getMessageTemplate()->getComplaintId()
            . ")";
    }


    /**
     * @param ReturnOperationNewInitialData $initialData
     * @param MessageRecipientEntity $recipient
     * @param MessageRecipientEntity $sender
     * @return string
     */
    public function composeBody(
        InitialData $initialData,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string
    {
        $params = $initialData->getMessageTemplate();
        $message = "There is a new complaint ("
            . $params->getComplaintId()
            . "). Reseller id: "
            . $initialData->getReseller()->getId();

        foreach ($params->toArray() as $key => $value) {
            $message .= "$key: $value";
        }

        return $message;
    }
}
