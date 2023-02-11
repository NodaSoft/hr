<?php

namespace App\ORM;

use Attribute;

#[Attribute(Attribute::TARGET_PROPERTY)]
class Column
{
    public function __construct(
        public readonly ColumnType $type = ColumnType::VARCHAR
    ){
    }
}