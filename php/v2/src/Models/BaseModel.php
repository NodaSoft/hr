<?php

namespace App\Models;

class BaseModel
{
    public int $id;
    public int $type;
    public string $name;

    public static function getById(int $id): ?self
    {
        return new static($id);
    }

    public function __construct(int $id)
    {
        $this->id = $id;
    }


    public function getFullName(): string
    {
        return sprintf("Name: %s, Id: %s", $this->name, $this->id);
    }
}
