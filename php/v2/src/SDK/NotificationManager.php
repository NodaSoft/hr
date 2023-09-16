<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Client;

use NW\WebService\References\Operations\Notification\Enum\Status;
use NW\WebService\References\Operations\Notification\Struct\Template;
use Symfony\Component\HttpFoundation\Exception\BadRequestException;

final class NotificationManager
{
    /**
     * @return non-empty-string[]
     * @throws BadRequestException
     */
    public static function send(int $sellerId, int $clientId, object $event, Status $status, Template $template): array
    {
        return ['success'];
    }
}
