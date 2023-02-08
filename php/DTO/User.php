<?php

declare(strict_types=1);

namespace NodaSoft\DTO;

final class User
{
    public function __construct(
        public int $id,
        public string $name,
        public string $lastName,
        public int $age,
        public ?string $from,
        public ?string $key,
    ) {
    }
}
