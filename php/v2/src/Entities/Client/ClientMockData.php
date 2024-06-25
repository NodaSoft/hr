<?php

namespace Nodasoft\Testapp\Entities;

use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Interfaces\MockDataInterface;

class ClientMockData implements MockDataInterface
{
    public static function get(): array
    {
        return [
            [
                'id' => 1,
                'name' => 'client 1',
                'email' => 'client1email@noda.soft',
                'mobile' => '+374855454545',
                'seller_id' => 5,
                'type' => ContactorType::TYPE_CUSTOMER->value,
            ],
            [
                'id' => 2,
                'name' => 'client 2',
                'email' => null,
                'mobile' => '+37454554545',
                'seller_id' => 7,
                'type' => ContactorType::TYPE_CUSTOMER->value,
            ],
            [
                'id' => 3,
                'name' => 'client 3',
                'email' => 'client3email@noda.soft',
                'mobile' => null,
                'seller_id' => 2,
                'type' => ContactorType::TYPE_CUSTOMER->value,
            ],
        ];
    }
}