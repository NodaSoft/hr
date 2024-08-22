<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Email;

class EmailMessage
{
    private string $emailFrom;
    private string $emailTo;
    private string $subject;
    private string $body;
    private int $resellerId;
    private string $eventType;
    private ?int $clientId;
    private ?int $newStatus;

    public function __construct(
        string $emailFrom,
        string $emailTo,
        string $subjectTemplate,
        string $bodyTemplate,
        int $resellerId,
        string $eventType,
        array $templateData,
        ?int $clientId = null,
        ?int $newStatus = null
    ) {
        $this->emailFrom = $emailFrom;
        $this->emailTo = $emailTo;
        $this->subject = __($subjectTemplate, $templateData, $resellerId);
        $this->body = __($bodyTemplate, $templateData, $resellerId);
        $this->resellerId = $resellerId;
        $this->eventType = $eventType;
        $this->clientId = $clientId;
        $this->newStatus = $newStatus;
    }
    public function send(): void
    {
        $messageData = [
            [
                'emailFrom' => $this->emailFrom,
                'emailTo' => $this->emailTo,
                'subject' => $this->subject,
                'message' => $this->body,
            ],
        ];

        MessagesClient::sendMessage($messageData, $this->resellerId, $this->clientId, $this->eventType, $this->newStatus);
    }
}