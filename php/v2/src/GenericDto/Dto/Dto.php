<?php

namespace NodaSoft\GenericDto\Dto;

interface Dto
{
    public function isValid(): bool;

    public function toArray(): array;
}
