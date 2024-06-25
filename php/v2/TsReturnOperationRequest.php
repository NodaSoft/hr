<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperationRequest
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
        $this->resellerId = (int)$data['resellerId'];
        $this->notificationType = (int)$data['notificationType'];
        $this->clientId = (int)$data['clientId'];
        $this->creatorId = (int)$data['creatorId'];
        $this->expertId = (int)$data['expertId'];
        $this->complaintId = (int)$data['complaintId'];
        $this->complaintNumber = (string)$data['complaintNumber'];
        $this->consumptionId = (int)$data['consumptionId'];
        $this->consumptionNumber = (string)$data['consumptionNumber'];
        $this->agreementNumber = (string)$data['agreementNumber'];
        $this->date = (string)$data['date'];
        $this->differences = $data['differences'] ?? [];
    }
}
