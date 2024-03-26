<?php

namespace src\Operation\Domain\Enum;

enum PositionStatus: int
{
    const COMPLETED = 0;
    const PENDING = 1;
    const REJECTED = 2;

    const STATUSES = [
        self::COMPLETED => 'completed',
        self::PENDING => 'pending',
        self::REJECTED => 'rejected',
    ];
}
