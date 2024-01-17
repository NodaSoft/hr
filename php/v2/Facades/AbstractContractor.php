<?php


namespace App\v2\Facades;

class AbstractContractor
{
    public const TYPE_CUSTOMER = 0;
    public int $id;
    public int $type;
    public string $name;

    /**
     * @param  int  $id
     * @return static
     */
    public static function getById(int $id): static
    {
        return new static($id); // fakes the getById method
    }

    /**
     * @return string
     */
    public function getFullName(): string
    {
        return $this->name.' '.$this->id;
    }
}
