<?php

namespace NodaSoft\GenericDto\Dto;

class ReturnOperationNewMessageContentList implements Dto
{
    /** @var int */
    private $COMPLAINT_ID;

    /** @var string */
    private $COMPLAINT_NUMBER;

    /** @var int */
    private $CREATOR_ID;

    /** @var string */
    private $CREATOR_NAME;

    /** @var int */
    private $EXPERT_ID;

    /** @var string */
    private $EXPERT_NAME;

    /** @var int */
    private $CLIENT_ID;

    /** @var string */
    private $CLIENT_NAME;

    /** @var int */
    private $CONSUMPTION_ID;

    /** @var string */
    private $CONSUMPTION_NUMBER;

    /** @var string */
    private $AGREEMENT_NUMBER;

    /** @var string */
    private $DATE;

    /**
     * @return array<string, mixed>
     */
    public function toArray(): array
    {
        $array = [];
        foreach ($this as $key => $value) {
            $array[$key] = $value;
        }
        return $array;
    }

    public function isValid(): bool
    {
        return empty($this->getEmptyKeys());
    }

    /**
     * @return string[]
     */
    public function getEmptyKeys(): array
    {
        $emptyKeys = [];
        foreach ($this as $key => $value) {
            if (empty($value)) {
                $emptyKeys[] = $key;
            }
        }
        return $emptyKeys;
    }

    public function getComplaintId(): int
    {
        return $this->COMPLAINT_ID;
    }

    public function setComplaintId(int $complaintId): void
    {
        $this->COMPLAINT_ID = $complaintId;
    }

    public function getComplaintNumber(): string
    {
        return $this->COMPLAINT_NUMBER;
    }

    public function setComplaintNumber(string $complaintNumber): void
    {
        $this->COMPLAINT_NUMBER = $complaintNumber;
    }

    public function getCreatorId(): int
    {
        return $this->CREATOR_ID;
    }

    public function setCreatorId(int $creatorId): void
    {
        $this->CREATOR_ID = $creatorId;
    }

    public function getCreatorName(): string
    {
        return $this->CREATOR_NAME;
    }

    public function setCreatorName(string $creatorName): void
    {
        $this->CREATOR_NAME = $creatorName;
    }

    public function getExpertId(): int
    {
        return $this->EXPERT_ID;
    }

    public function setExpertId(int $expertId): void
    {
        $this->EXPERT_ID = $expertId;
    }

    public function getExpertName(): string
    {
        return $this->EXPERT_NAME;
    }

    public function setExpertName(string $expertName): void
    {
        $this->EXPERT_NAME = $expertName;
    }

    public function getClientId(): int
    {
        return $this->CLIENT_ID;
    }

    public function setClientId(int $clientId): void
    {
        $this->CLIENT_ID = $clientId;
    }

    public function getClientName(): string
    {
        return $this->CLIENT_NAME;
    }

    public function setClientName(string $clientName): void
    {
        $this->CLIENT_NAME = $clientName;
    }

    public function getConsumptionId(): int
    {
        return $this->CONSUMPTION_ID;
    }

    public function setConsumptionId(int $consumptionId): void
    {
        $this->CONSUMPTION_ID = $consumptionId;
    }

    public function getConsumptionNumber(): string
    {
        return $this->CONSUMPTION_NUMBER;
    }

    public function setConsumptionNumber(string $consumptionNumber): void
    {
        $this->CONSUMPTION_NUMBER = $consumptionNumber;
    }

    public function getAgreementNumber(): string
    {
        return $this->AGREEMENT_NUMBER;
    }

    public function setAgreementNumber(string $agreementNumber): void
    {
        $this->AGREEMENT_NUMBER = $agreementNumber;
    }

    public function getDate(): string
    {
        return $this->DATE;
    }

    public function setDate(string $date): void
    {
        $this->DATE = $date;
    }
}
