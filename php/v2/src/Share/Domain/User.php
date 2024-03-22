<?php

namespace Src\Share\Domain;

class User
{
    public int $id;
    public int $type;
    public string $name;
    public ?int $sellerId = null;
    public ?string $email = null;
    public ?string $mobile = null;

    public static function getById(int $id): self
    {
        return new self($id);
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }

}