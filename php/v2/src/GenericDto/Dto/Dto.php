<?php

namespace NodaSoft\GenericDto\Dto;

interface Dto
{
    public function isValid(): bool;

    /**
     * @return array<string, mixed>
     */
    public function toArray(): array;
}
