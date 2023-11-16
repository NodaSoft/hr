<?php

namespace NodaSoft\Operation\Params;

use NodaSoft\Request\Request;

interface Params
{
    public function setRequest(Request $request): void;

    public function isValid(): bool;

    /**
     * @return array<string, mixed>
     */
    public function toArray(): array;
}
