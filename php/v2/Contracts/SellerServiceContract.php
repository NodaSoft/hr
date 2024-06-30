<?php

namespace NW\WebService\References\Operations\Notification\Contracts;

use NW\WebService\References\Operations\Notification\Seller;

interface SellerServiceContract
{
    public function getById(int $id): ?Seller;
}
