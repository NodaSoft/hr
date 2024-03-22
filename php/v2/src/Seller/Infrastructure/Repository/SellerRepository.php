<?php

namespace src\Seller\Infrastructure\Repository;

use Src\Seller\Domain\Entity\Seller;
use Src\Seller\Domain\Repository\SellerRepositoryInterface;

class SellerRepository implements SellerRepositoryInterface

{
    public function getById(int $id): Seller
    {
        // TODO: Implement getById() method.
    }
}