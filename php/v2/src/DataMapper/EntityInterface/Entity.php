<?php

namespace NodaSoft\DataMapper\EntityInterface;

interface Entity
{
    /**
     * @return array<string, mixed>
     */
    public function toArray(): array;

    public function getId(): int;

    public function setId(int $id): void;

    public function setName(string $name): void;

    public function getName(): string;
}
