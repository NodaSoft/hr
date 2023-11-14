<?php

namespace NodaSoft\Message\Template;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationStatusChangedInitialData;

class ComplaintStatusChanged implements Template
{
    /**
     * @param ReturnOperationStatusChangedInitialData $data
     * @param MessageRecipientEntity $recipient
     * @param MessageRecipientEntity $sender
     * @return string
     */
    public function composeSubject(
        InitialData $data,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string
    {
        return "Complaint status has been changed ("
            . $data->getMessageTemplate()->getComplaintId()
            . ")";
    }

    /**
     * @param ReturnOperationStatusChangedInitialData $data
     * @param MessageRecipientEntity $recipient
     * @param MessageRecipientEntity $sender
     * @return string
     */
    public function composeBody(
        InitialData $data,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string {
        $params = $data->getMessageTemplate();
        $message = "Hello, " . $recipient->getFullName()
            . ". Be informed that complaint №" . $params->getComplaintId()
            . " status has been changed from " . $data->getPreviousStatusName()
            . " to " . $data->getCurrentStatusName()
            . "). Reseller id: " . $data->getReseller()->getId();

        foreach ($params->toArray() as $key => $value) {
            $message .= "$key: $value";
        }

        return $message;
    }
}
