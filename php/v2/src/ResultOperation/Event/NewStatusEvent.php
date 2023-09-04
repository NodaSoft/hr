<?php

declare(strict_types=1);

namespace ResultOperation\Event;

use ResultOperation\DTO\NotificationTemplate;
use ResultOperation\Entity\Contractor;

class NewStatusEvent extends AbstractStatusEvent
{
    public function __construct(
        protected readonly Contractor $client,
        protected readonly NotificationTemplate $template,
    ) {
    }
}
