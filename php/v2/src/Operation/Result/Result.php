<?php

namespace NodaSoft\Operation\Result;

interface Result
{
    /**
     * @return array<string, mixed>
     */
    public function toArray(): array;
}
