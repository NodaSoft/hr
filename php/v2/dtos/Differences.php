<?php

namespace NW\WebService\References\Operations\Notification;

final readonly class Differences
{
    public function __construct(
        public int $from,
        public int $to,
    )
    {
    }
}