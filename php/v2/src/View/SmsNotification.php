<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Struct;

final class SmsNotification
{
    /** @phpstan-ignore-next-line */
    private bool $isSent;

    /** @phpstan-ignore-next-line */
    private string $message;

    public function __construct(bool $isSent = false, string $message = '')
    {
        $this->isSent  = $isSent;
        $this->message = $message;
    }
}
