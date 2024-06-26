<?php

namespace Nodasoft\Testapp\Interfaces;

use Nodasoft\Testapp\Enums\ContactorType;

interface ContactorInterface
{
    const int TYPE_CUSTOMER = 0;

    public function getId(): int;

    public function getType(): ContactorType;

    public function getFullName(): string;
}