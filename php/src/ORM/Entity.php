<?php

namespace App\ORM;

use Attribute;

#[Attribute(Attribute::TARGET_CLASS)]
class Entity
{
    public function __construct(
        public readonly string $repository,
        public readonly ?string $table = null,
    ){
    }
}