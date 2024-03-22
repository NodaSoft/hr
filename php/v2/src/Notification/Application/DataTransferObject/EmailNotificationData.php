<?php

namespace Src\Notification\Application\DataTransferObject;

class EmailNotificationData
{
    public string $emailFrom;
    public array $emails;
    public string $subject;
    public string $message;
    public int $resellerId;

    public function __construct(string $emailFrom, array $emails, string $subject, string $message, int $resellerId)
    {
        $this->emailFrom = $emailFrom;
        $this->emails = $emails;
        $this->subject = $subject;
        $this->message = $message;
        $this->resellerId = $resellerId;
    }

    public static function fromArray(array $data): self
    {
        return new self(
            $data['emailFrom'],
            $data['emails'],
            $data['subject'],
            $data['message'],
            $data['resellerId']
        );
    }

    public function toArray(): array
    {
        return [
            'emailFrom' => $this->emailFrom,
            'emails' => $this->emails,
            'subject' => $this->subject,
            'message' => $this->message,
            'resellerId' => $this->resellerId,
        ];
    }
}