<?php

namespace Nodasoft\Testapp\Entities\Seller;

use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Interfaces\MockDataInterface;

class SellerMockData implements MockDataInterface
{
    public static function get(): array
    {
        return [
            [
                'id' => 1,
                'name' => 'seller 1',
                'email' => 'seller1email@noda.soft',
                'type' => ContactorType::TYPE_SELLER->value,
            ],
            [
                'id' => 2,
                'name' => 'seller 2',
                'email' => 'seller2email@noda.soft',
                'type' => ContactorType::TYPE_SELLER->value,
            ],
            [
                'id' => 3,
                'name' => 'seller 3',
                'email' => 'seller3email@noda.soft',
                'type' => ContactorType::TYPE_SELLER->value,
            ],
        ];
    }
}