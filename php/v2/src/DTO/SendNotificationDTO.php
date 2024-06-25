<?php

namespace Nodasoft\Testapp\DTO;

use Nodasoft\Testapp\Enums\NotificationType;

class SendNotificationDTO
{
    public function __construct(
        public int                   $resellerId,
        public NotificationType      $notificationType,
        public int                   $clientId,
        public int                   $creatorId,
        public int                   $expertId,
        public int                   $complaintId,
        public string                $complaintNumber,
        public int                   $consumptionId,
        public string                $consumptionNumber,
        public string                $agreementNumber,
        public string                $date,
        public ?MessageDifferenceDto $differences = null,
    )
    {
    }
}