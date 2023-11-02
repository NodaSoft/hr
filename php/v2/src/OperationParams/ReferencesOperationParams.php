<?php

namespace NodaSoft\OperationParams;

use NodaSoft\Request\Request;

interface ReferencesOperationParams
{
    public function setRequest(Request $request): void;

    public function isValid(): bool;

    public function toArray(): array;
}
