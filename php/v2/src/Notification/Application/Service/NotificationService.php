<?php

namespace Src\Notification\Application\Service;

use Src\Notification\Application\DataTransferObject\EmailNotificationData;
use Src\Notification\Application\DataTransferObject\SmsNotificationData;
use src\Notification\Domain\Enum\NotificationStatus;

readonly class NotificationService
{
    private EmailService $emailService;
    private SmsService $smsService;

    public function __construct()
    {
        $this->emailService = new EmailService();
        $this->smsService = new SmsService();
    }

    public function sendEmailNotification(EmailNotificationData $data): void
    {
        if (!empty($data->emailFrom) && !empty($data->emails)) {
            foreach ($data->emails as $email) {
                $this->emailService->send([
                    'emailFrom' => $data->emailFrom,
                    'emailTo' => $email,
                    'subject' => __('complaintClientEmailSubject', $data->message, $data->resellerId),
                    'message' => __('complaintClientEmailBody', $data->message, $data->resellerId),
                ], $data->resellerId, NotificationStatus::CHANGE_RETURN_STATUS);
            }
        }
    }

    public function sendSmsNotification(SmsNotificationData $data): void
    {
        if (!empty($data->phoneNumber)) {
            $this->smsService->send($data);
        }
    }

}