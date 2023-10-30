<?php

namespace NodaSoft\OperationParams;

use NodaSoft\Request\Request;

interface OperationParams
{
    public function setRequest(Request $request): void;

    public function isValid(): bool;

    public function toArray(): array;
}
