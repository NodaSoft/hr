<?php

namespace NW\WebService\References\Operations\Notification\DTOs;

class StatusEnum
{
    public const COMPLETED = 0;
    public const PENDING = 1;
    public const REJECTED = 2;

    public const STATUS_MAP = [
        self::COMPLETED => 'Completed',
        self::PENDING => 'Pending',
        self::REJECTED => 'Rejected',
    ];
}
