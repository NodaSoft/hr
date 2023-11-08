<?php

namespace NodaSoft\DataMapper\EntityTrait;

trait EmailEntity
{
    use Entity;

    /** @var null|string */
    private $email;

    public function setEmail(?string $email): void
    {
        $this->email = $email;
    }

    public function hasEmail(): bool
    {
        return is_null($this->getEmail());
    }

    public function getEmail(): ?string
    {
        return $this->email;
    }

    public function toArray(): array
    {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'email' => $this->email,
        ];
    }
}
