<?php

namespace Src\Operation\Application\DataTransferObject;

class OperationData
{
    public int $resellerId;
    public int $clientId;
    public int $notificationType;
    public int $creatorId;
    public int $expertId;
    public int $complaintId;
    public string $complaintNumber;
    public int $consumptionId;
    public string $consumptionNumber;
    public string $agreementNumber;
    public string $date;
    public DifferencesData $differences;

    /**
     * @param array $data
     * @return self
     */
    public static function fromArray(array $data): self
    {
        $dto = new self();
        $dto->resellerId = $data['resellerId'];
        $dto->clientId = $data['clientId'];
        $dto->notificationType = $data['notificationType'];
        $dto->creatorId = $data['creatorId'];
        $dto->expertId = $data['expertId'];
        $dto->complaintId = $data['complaintId'];
        $dto->complaintNumber = $data['complaintNumber'];
        $dto->consumptionId = $data['consumptionId'];
        $dto->consumptionNumber = $data['consumptionNumber'];
        $dto->agreementNumber = $data['agreementNumber'];
        $dto->date = $data['date'];
        $dto->differences = DifferencesData::fromArray($data['differences']);

        return $dto;
    }
}