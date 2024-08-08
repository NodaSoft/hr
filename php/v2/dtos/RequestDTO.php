<?php

namespace NW\WebService\References\Operations\Notification;

final readonly class RequestDTO
{
    public function __construct(
        public int $resellerId,
        public int $clientId,
        public int $creatorId,
        public int $expertId,
        public Differences $differences,
        public ?NotificationType $notificationType,
        public int $complaintId,
        public int $consumptionId,
        public string $complaintNumber,
        public string $consumptionNumber,
        public string $agreementNumber,
        public string $date
    ) {}
}