<?php

namespace NodaSoft\Message\Template;

use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\TsReturnInitialData;

class ComplaintStatus implements Template
{
    /**
     * @param TsReturnInitialData $initialData
     * @return string
     */
    public function composeSubject(InitialData $initialData): string
    {
        //todo: handle template logic

        return "Complaint claim ("
            . $initialData->getMessageTemplate()->getDate()
            . ")";
    }

    /**
     * @param TsReturnInitialData $initialData
     * @return string
     */
    public function composeBody(InitialData $initialData): string
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
}
