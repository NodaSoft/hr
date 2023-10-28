<?php

namespace NodaSoft\DataMapper\Entity;

class Reseller implements Entity
{
    /** @var int */
    private $id;

    /** @var string */
    private $name;

    public function setId(int $id): void
    {
        $this->id = $id;
    }

    public function getId(): int
    {
        return $this->id;
    }

    public function setName(string $name): void
    {
        $this->name = $name;
    }

    public function getName(): string
    {
        return $this->name;
    }

    public function getFullName(): string
    {
        return $this->getName() . ' ' . $this->getId();
    }
}
