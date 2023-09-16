<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Enum;

enum Status: string
{
    case COMPLETED = 'Completed';

    case PENDING = 'Pending';

    case REJECTED = 'Rejected';
}
