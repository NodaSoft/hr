<?php

declare(strict_types=1);

namespace App\Domain;

final class User
{
    public function __construct(
        public string $firstName,
        public string $lastName,
        public string $location,
        public int $age,
        public array $settings = [],
        public ?int $id = null,
    )
    {
    }
}
