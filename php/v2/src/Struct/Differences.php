<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Struct;

use NW\WebService\References\Operations\Notification\Enum\Status;

final class Differences
{
    public function __construct(public Status $from, public Status $to)
    {
    }
}
