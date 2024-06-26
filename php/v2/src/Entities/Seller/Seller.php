<?php

namespace Nodasoft\Testapp\Entities\Seller;

use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Interfaces\ContactorInterface;

readonly class Seller implements ContactorInterface
{
    public function __construct(
        private int           $id,
        private string        $name,
        private string        $email,
        private ContactorType $type
    )
    {
    }

    public function getType(): ContactorType
    {
        return $this->type;
    }

    public function getId(): int
    {
        return $this->id;
    }

    public function getName(): string
    {
        return $this->name;
    }

    public function getEmail(): string
    {
        return $this->email;
    }

    public function getFullName(): string
    {
        return $this->name;
    }
}