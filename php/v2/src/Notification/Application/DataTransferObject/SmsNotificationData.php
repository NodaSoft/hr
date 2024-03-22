<?php

namespace Src\Notification\Application\DataTransferObject;

class SmsNotificationData
{
    public string $phoneNumber;
    public string $message;
    public int $resellerId;
    public int $clientId;
    public int $status;

    public function __construct(string $phoneNumber, string $message, int $resellerId, int $clientId, int $status)
    {
        $this->phoneNumber = $phoneNumber;
        $this->message = $message;
        $this->resellerId = $resellerId;
        $this->clientId = $clientId;
        $this->status = $status;
    }

    public static function fromArray(array $data): self
    {
        return new self(
            $data['phoneNumber'],
            $data['message'],
            $data['resellerId'],
            $data['clientId'],
            $data['status']
        );
    }

    public function toArray(): array
    {
        return [
            'phoneNumber' => $this->phoneNumber,
            'message' => $this->message,
            'resellerId' => $this->resellerId,
            'clientId' => $this->clientId,
            'status' => $this->status,
        ];
    }
}