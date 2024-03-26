<?php

namespace Src\Operation\Infrastructure\Adapters;

use Src\Notification\Infrastructure\API\NotificationApi;

readonly class NotificationAdapter
{
    private NotificationApi $notificationApi;

    public function __construct()
    {
        $this->notificationApi = new NotificationApi();
    }

    public function sendEmailNotification(array $emailData): array
    {
        return $this->notificationApi->sendEmailNotification($emailData);
    }

    public function sendSmsNotification(array $smsData): array
    {
        return $this->notificationApi->sendSmsNotification($smsData);
    }

}