<?php

namespace NodaSoft\DataMapper\EntityInterface;

interface EmailEntity
{
    public function setEmail(string $email): void;

    public function hasEmail(): bool;

    public function getEmail(): ?string;
}
