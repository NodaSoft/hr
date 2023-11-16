<?php

namespace NodaSoft\Messenger;

interface Recipient
{
    public function getFullName(): string;

    public function setEmail(string $email): void;

    public function hasEmail(): bool;

    public function getEmail(): ?string;

    public function getCellphone(): ?int;

    public function setCellphone(int $number): void;

    public function hasCellphone(): bool;

    public function toArray(): array;
}
