<?php

namespace App\v2\Requests;
use Exception;

class NotificationRequest
{
    private array $data;
    public function __construct(array $request)
    {
        $this->data = $request['data'];
    }

    /**
     * @throws Exception
     */
    public function validated() : array
    {
        return [
            'clientId' => $this->data['clientId'] ?? throw new Exception('Client id not found', 422),
            'resellerId' => $this->data['resellerId'] ?? throw new Exception('Reseller id not found', 422),
            'notificationType' => $this->data['notificationType'] ?? throw new Exception('Notification type not found', 422),
            'creatorId' => $this->data['creatorId'] ?? throw new Exception('Creator id not found', 422),
            'expertId' => $this->data['expertId'] ?? throw new Exception('Expert id not found', 422),
            'differences' => $this->data['differences'] ?? throw new Exception('Differences not found', 422),
            'date' => $this->data['date'] ?? throw new Exception('Date not found', 422),
            'consumptionId' => $this->data['consumptionId'] ?? throw new Exception('Consumption ID not found', 422),
            'consumptionNumber' => $this->data['consumptionNumber'] ?? throw new Exception('Consumption number not found', 422),
            'agreementNumber' => $this->data['agreementNumber'] ?? throw new Exception('Agreement number not found', 422),
            'complaintId' => $this->data['complaintId'] ?? throw new Exception('Complaint ID not found', 422),
            'complaintNumber' => $this->data['complaint'] ?? throw new Exception('Complaint not found', 422),
        ];
    }
}