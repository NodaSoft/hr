<?php

namespace NodaSoft\ReferencesOperation\Params;

use NodaSoft\Request\Request;

interface ReferencesOperationParams
{
    public function setRequest(Request $request): void;

    public function isValid(): bool;

    public function toArray(): array;
}
