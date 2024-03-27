<?php

namespace App\Models;

class Contractor extends BaseModel
{
    public Seller $seller;

    public ?string $email;
    public ?string $mobile;
    public function getSeller(): Seller
    {
        return $this->seller;
    }

    public function setSeller(Seller $seller): Contractor
    {
        $this->seller = $seller;
        return $this;
    }
}
