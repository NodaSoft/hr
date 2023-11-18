<?php

namespace NodaSoft\DataMapper\EntityTrait;

trait MessageRecipientEntity
{
    use Entity;

    /** @var ?string */
    private $email = null;

    /** @var ?int */
    private $cellphone = null;

    public function setEmail(string $email): void
    {
        $this->email = $email;
    }

    public function hasEmail(): bool
    {
        return ! is_null($this->email);
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

    /**
     * @return array<string, mixed>
     */
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
