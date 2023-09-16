<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Client;

use NW\WebService\References\Operations\Notification\Enum\Status;
use NW\WebService\References\Operations\Notification\Struct\Email;

interface Client
{
    public static function sendMessage(Email $email, int $sellerId, object $event, ?int $clientId = null, ?Status $status = null): void;
}
