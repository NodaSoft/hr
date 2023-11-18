<?php

namespace NodaSoft\GenericDto\Dto;

interface Dto
{
    /**
     * @return array<string, int|string>
     */
    public function toArray(): array;
}
