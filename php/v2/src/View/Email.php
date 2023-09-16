<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Struct;

final class Email
{
    public function __construct(
        public readonly string $from,
        public readonly string $to,
        public readonly string $subject,
        public readonly string $body
    ) {
    }
}
