<?php

namespace Nodasoft\Testapp\DTO;

class MessageDifferenceDto
{
    public function __construct(
        public int $from,
        public int $to
    )
    {
    }
}