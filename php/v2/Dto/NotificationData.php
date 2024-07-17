<?php

namespace NW\WebService\References\Operations\Notification\Dto;

use NW\WebService\References\Operations\Notification\Notification\Enums\NotificationTypeEnum;
use NW\WebService\References\Operations\Notification\Status;

class NotificationData
{
    public function __construct(
        public int                  $resellerId,
        public NotificationTypeEnum $notificationType,
        public int                  $clientId,
        public int                  $creatorId,
        public int                  $expertId,
        public int                  $complaintId,
        public string               $complaintNumber,
        public int                  $consumptionId,
        public string               $consumptionNumber,
        public string               $agreementNumber,
        public string               $date,
        public array|string         $differences
    )
    {
        $this->differences = $this->getDifferences();
    }

    private function getDifferences(): string
    {
        if ($this->notificationType === NotificationTypeEnum::NEW) {
            return 'New position added';
        }

        if ($this->notificationType === NotificationTypeEnum::CHANGE && !empty($this->differences)) {
            return sprintf(
                'Status changed from %s to %s',
                Status::getName((int)$this->differences['from']),
                Status::getName((int)$this->differences['to'])
            );
        }

        return '';
    }
}