<?php

declare(strict_types=1);

namespace NodaSoft\DTO;

final class NewUser
{
    public function __construct(
        public string $name,
        public string $lastName,
        public int $age,
    ) {
    }
}
