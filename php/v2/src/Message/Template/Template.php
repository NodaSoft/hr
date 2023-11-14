<?php

namespace NodaSoft\Message\Template;

use NodaSoft\DataMapper\EntityInterface\MessageRecipientEntity;
use NodaSoft\ReferencesOperation\InitialData\InitialData;

interface Template
{
    public function composeSubject(
        InitialData $data,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string;

    public function composeBody(
        InitialData $data,
        MessageRecipientEntity $recipient,
        MessageRecipientEntity $sender
    ): string;
}
