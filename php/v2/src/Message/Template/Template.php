<?php

namespace NodaSoft\Message\Template;

use NodaSoft\ReferencesOperation\InitialData\InitialData;

interface Template
{
    public function composeSubject(InitialData $initialData): string;

    public function composeBody(InitialData $initialData): string;
}
