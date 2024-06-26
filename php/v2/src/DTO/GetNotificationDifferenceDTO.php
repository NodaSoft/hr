<?php

namespace Nodasoft\Testapp\DTO;

use Nodasoft\Testapp\Enums\NotificationType;

class GetNotificationDifferenceDTO
{
    public function __construct(
        public NotificationType      $notificationType,
        public int                   $resellerId,
        public ?MessageDifferenceDto $messageDifference = null
    )
    {
    }
}