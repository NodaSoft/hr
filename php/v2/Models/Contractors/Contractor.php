<?php


namespace App\v2\Models\Contractors;

use App\v2\Facades\AbstractContractor;
/**
 * @property Seller $Seller
 */
class Contractor extends AbstractContractor
{
    public const int TYPE_CUSTOMER = 0;
    public int $id;
    public int $type;
    public string $name;
    public string $email;
    public string $mobile;
}
