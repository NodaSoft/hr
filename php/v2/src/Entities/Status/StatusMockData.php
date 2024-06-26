<?php

namespace Nodasoft\Testapp\Entities;


use Nodasoft\Testapp\Interfaces\MockDataInterface;

class StatusMockData implements MockDataInterface
{
    public static function get(): array
    {
        return [
            [
                'id' => 1,
                'name' => 'Completed'
            ],
            [
                'id' => 2,
                'name' => 'Pending'
            ],
            [
                'id' => 3,
                'name' => 'Rejected'
            ]
        ];
    }
}