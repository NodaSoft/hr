<?php

namespace  Nodasoft\Testapp\Repositories\Seller;

use Exception;
use Nodasoft\Testapp\Entities\Seller\Seller;
use Nodasoft\Testapp\Entities\Seller\SellerMockData;
use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Traits\CanGetByKey;

class SellerRepository implements SellerRepositoryInterface
{
    use CanGetByKey;

    /**
     * @throws Exception
     */
    public function getById(int $id): Seller
    {
        $record = $this->getByKeyOrThrow(SellerMockData::get(), $id);

        return new Seller(
            $record['id'],
            $record['name'],
            $record['email'],
            ContactorType::tryFrom($record['type'])
        );
    }
}