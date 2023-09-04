<?php

declare(strict_types=1);

namespace ResultOperation\Enum;

enum Status: int
{
    case Completed = 0;
    case Pending = 1;
    case Rejected = 2;
}
