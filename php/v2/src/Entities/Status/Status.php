<?php

namespace Nodasoft\Testapp\Entities\Status;

readonly class Status
{
    public function __construct(
        private int    $id,
        private string $name
    )
    {
    }


    public function getName(): string
    {
        return $this->name;
    }

    public function getId(): int
    {
        return $this->id;
    }
}