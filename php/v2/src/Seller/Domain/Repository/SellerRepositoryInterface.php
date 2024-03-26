<?php

namespace Src\Seller\Domain\Repository;

use Src\Seller\Domain\Entity\Seller;

interface SellerRepositoryInterface
{
    public function getById(int $id): Seller;

}