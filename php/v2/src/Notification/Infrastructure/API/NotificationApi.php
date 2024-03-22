<?php

namespace Src\Notification\Infrastructure\API;

use Src\Notification\Application\DataTransferObject\EmailNotificationData;
use Src\Notification\Application\DataTransferObject\SmsNotificationData;
use Src\Notification\Application\Service\NotificationService;

class NotificationApi
{
    private NotificationService $service;

    public function __construct()
    {
        $this->service = new NotificationService();

    }

    public function sendEmailNotification(array $emailData): void
    {
        $data = EmailNotificationData::fromArray($emailData);
        $this->service->sendEmailNotification($data);
    }

    public function sendSmsNotification(array $smsData): void
    {
        $data = SmsNotificationData::fromArray($smsData);
        $this->service->sendSmsNotification($data);
    }

}