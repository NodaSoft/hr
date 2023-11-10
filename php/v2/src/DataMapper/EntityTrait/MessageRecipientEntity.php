<?php

namespace NodaSoft\DataMapper\EntityTrait;

trait MessageRecipientEntity
{
    use Entity;

    /** @var null|string */
    private $email;

    /** @var int */
    private $cellphone;

    public function setEmail(?string $email): void
    {
        $this->email = $email;
    }

    public function hasEmail(): bool
    {
        return ! is_null($this->getEmail());
    }

    public function getEmail(): ?string
    {
        return $this->email;
    }

    public function getCellphone(): ?int
    {
        return $this->cellphone ?? null;
    }

    public function setCellphone(int $number): void
    {
        $this->cellphone = $number;
    }

    public function hasCellphone(): bool
    {
        return ! is_null($this->getCellphone());
    }

    public function toArray(): array
    {
        return [
            'id' => $this->getId(),
            'name' => $this->getName(),
            'email' => $this->getEmail(),
            'cellphone' => $this->getCellphone()
        ];
    }
}
