<?php

namespace NodaSoft\Message\Template;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationStatusChangedInitialData;

class ComplaintStatusChanged implements Template
{
    /**
     * @param ReturnOperationStatusChangedInitialData $initialData
     * @param MessageRecipientEntity $recipient
     * @param MessageRecipientEntity $sender
     * @return string
     */
    public function composeSubject(
        InitialData $initialData,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string
    {
        return "Complaint status has been changed ("
            . $initialData->getMessageTemplate()->getComplaintId()
            . ")";
    }

    /**
     * @param ReturnOperationStatusChangedInitialData $initialData
     * @param MessageRecipientEntity $recipient
     * @param MessageRecipientEntity $sender
     * @return string
     */
    public function composeBody(
        InitialData $initialData,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string {
        $params = $initialData->getMessageTemplate();
        $message = "Complaint status has been changed ("
            . $params->getComplaintId()
            . "). Reseller id: "
            . $initialData->getReseller()->getId();

        foreach ($params->toArray() as $key => $value) {
            $message .= "$key: $value";
        }

        return $message;
    }
}
