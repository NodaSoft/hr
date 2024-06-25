<?php

namespace Nodasoft\Testapp\Repositories\Seller;


use Nodasoft\Testapp\Entities\Seller\Seller;

interface SellerRepositoryInterface
{
    public function getById(int $id): Seller;
}