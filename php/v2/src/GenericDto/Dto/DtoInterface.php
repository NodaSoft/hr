<?php

namespace NodaSoft\GenericDto\Dto;

interface DtoInterface
{
    public function isValid(): bool;

    public function toArray(): array;
}
