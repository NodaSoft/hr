<?php

namespace NW\WebService\References\Operations\Notification;

class NotificationDTO
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    private ?int $resellerId;
    private ?int $type;
    private ?int $clientId;
    private ?int $creatorId;
    private ?int $expertId;
    private ?int $complaintId;
    private ?string $complaintNumber;
    private ?int $consumptionId;
    private ?string $consumptionNumber;
    private ?string $agreementNumber;
    private ?string $date;
    private ?array $differences = [];

    /**
     * @return int|null
     */
    public function getResellerId(): ?int
    {
        return $this->resellerId;
    }

    /**
     * @param int|null $resellerId
     * @return NotificationDTO
     */
    public function setResellerId(?int $resellerId): NotificationDTO
    {
        $this->resellerId = $resellerId;
        return $this;
    }

    /**
     * @return int|null
     */
    public function getType(): ?int
    {
        return $this->type;
    }

    /**
     * @param int|null $type
     * @return NotificationDTO
     */
    public function setType(?int $type): NotificationDTO
    {
        $this->type = $type;
        return $this;
    }

    /**
     * @return int|null
     */
    public function getClientId(): ?int
    {
        return $this->clientId;
    }

    /**
     * @param int|null $clientId
     * @return NotificationDTO
     */
    public function setClientId(?int $clientId): NotificationDTO
    {
        $this->clientId = $clientId;
        return $this;
    }

    /**
     * @return int|null
     */
    public function getCreatorId(): ?int
    {
        return $this->creatorId;
    }

    /**
     * @param int|null $creatorId
     * @return NotificationDTO
     */
    public function setCreatorId(?int $creatorId): NotificationDTO
    {
        $this->creatorId = $creatorId;
        return $this;
    }

    /**
     * @return int|null
     */
    public function getExpertId(): ?int
    {
        return $this->expertId;
    }

    /**
     * @param int|null $expertId
     * @return NotificationDTO
     */
    public function setExpertId(?int $expertId): NotificationDTO
    {
        $this->expertId = $expertId;
        return $this;
    }

    /**
     * @return int|null
     */
    public function getComplaintId(): ?int
    {
        return $this->complaintId;
    }

    /**
     * @param int|null $complaintId
     * @return NotificationDTO
     */
    public function setComplaintId(?int $complaintId): NotificationDTO
    {
        $this->complaintId = $complaintId;
        return $this;
    }

    /**
     * @return string|null
     */
    public function getComplaintNumber(): ?string
    {
        return $this->complaintNumber;
    }

    /**
     * @param string|null $complaintNumber
     * @return NotificationDTO
     */
    public function setComplaintNumber(?string $complaintNumber): NotificationDTO
    {
        $this->complaintNumber = $complaintNumber;
        return $this;
    }

    /**
     * @return int|null
     */
    public function getConsumptionId(): ?int
    {
        return $this->consumptionId;
    }

    /**
     * @param int|null $consumptionId
     * @return NotificationDTO
     */
    public function setConsumptionId(?int $consumptionId): NotificationDTO
    {
        $this->consumptionId = $consumptionId;
        return $this;
    }

    /**
     * @return string|null
     */
    public function getConsumptionNumber(): ?string
    {
        return $this->consumptionNumber;
    }

    /**
     * @param string|null $consumptionNumber
     * @return NotificationDTO
     */
    public function setConsumptionNumber(?string $consumptionNumber): NotificationDTO
    {
        $this->consumptionNumber = $consumptionNumber;
        return $this;
    }

    /**
     * @return string|null
     */
    public function getAgreementNumber(): ?string
    {
        return $this->agreementNumber;
    }

    /**
     * @param string|null $agreementNumber
     * @return NotificationDTO
     */
    public function setAgreementNumber(?string $agreementNumber): NotificationDTO
    {
        $this->agreementNumber = $agreementNumber;
        return $this;
    }

    /**
     * @return string|null
     */
    public function getDate(): ?string
    {
        return $this->date;
    }

    /**
     * @param string|null $date
     * @return NotificationDTO
     */
    public function setDate(?string $date): NotificationDTO
    {
        $this->date = $date;
        return $this;
    }

    /**
     * @return array|null
     */
    public function getDifferences(): ?array
    {
        return $this->differences;
    }

    /**
     * @param array|null $differences
     * @return NotificationDTO
     */
    public function setDifferences(?array $differences): NotificationDTO
    {
        $this->differences = $differences;
        return $this;
    }
}