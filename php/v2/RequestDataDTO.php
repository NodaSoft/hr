<?php

namespace NW\WebService\References\Operations\Notification;

class RequestDataDTO
{
    public int $resellerId;
    public int $notificationType;
    public int $clientId;
    public int $creatorId;
    public int $expertId;
    public int $complaintId;
    public string $complaintNumber;
    public int $consumptionId;
    public string $consumptionNumber;
    public string $agreementNumber;
    public string $date;
    public array $differences;

    public function __construct(array $data)
    {
        $this->resellerId = $data['resellerId'] ?? null;
        $this->notificationType = $data['notificationType'] ?? null;
        $this->clientId = $data['clientId'] ?? null;
        $this->creatorId = $data['creatorId'] ?? null;
        $this->expertId = $data['expertId'] ?? null;
        $this->complaintId = $data['complaintId'] ?? null;
        $this->complaintNumber = $data['complaintNumber'] ?? null;
        $this->consumptionId = $data['consumptionId'] ?? null;
        $this->consumptionNumber = $data['consumptionNumber'] ?? null;
        $this->agreementNumber = $data['agreementNumber'] ?? null;
        $this->date = $data['date'] ?? null;
        $this->differences = $data['differences'] ?? null;
    }
}