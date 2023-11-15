<?php

namespace NodaSoft\Messenger;

use NodaSoft\GenericDto\Dto\Dto;

interface Content
{
    public function composeMessageSubject(Dto $params): string;

    public function composeMessageBody(Dto $params): string;
}
