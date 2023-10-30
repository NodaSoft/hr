<?php

namespace NodaSoft\OperationParams;

use NodaSoft\Request\Request;

class TsReturnOperationParams implements OperationParams
{
    /** @var ?int */
    private $resellerId;

    /** @var ?int */
    private $clientId;

    /** @var ?int */
    private $creatorId;

    /** @var ?int */
    private $expertId;

    /** @var ?int */
    private $notificationType;

    /** @var ?int */
    private $complaintId;

    /** @var ?string */
    private $complaintNumber;

    /** @var ?int */
    private $consumptionId;

    /** @var ?string */
    private $consumptionNumber;

    /** @var ?string */
    private $agreementNumber;

    /** @var ?string */
    private $date;

    /** @var ?string */
    private $differences;

    public function setRequest(Request $request): void
    {
        foreach ($this as $key => $value) {
            $setter = 'set' . $key;
            if (method_exists($this, 'setter')) {
                $this->$setter($request->getData($key));
            }
        }
    }

    public function isValid(): bool
    {
        if (empty($resellerId)) {
            return false;
        }

        if (empty($notificationType)) {
            return false;
        }

        return true;
    }

    public function toArray(): array
    {
        $array = [];
        foreach ($this as $key => $value) {
            $array[$key] = $value;
        }
        return $array;
    }

    /**
     * @return int|null
     */
    public function getResellerId(): ?int
    {
        return $this->resellerId;
    }

    /**
     * @param int|null $resellerId
     */
    public function setResellerId(?int $resellerId): void
    {
        $this->resellerId = $resellerId;
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
     */
    public function setClientId(?int $clientId): void
    {
        $this->clientId = $clientId;
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
     */
    public function setCreatorId(?int $creatorId): void
    {
        $this->creatorId = $creatorId;
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
     */
    public function setExpertId(?int $expertId): void
    {
        $this->expertId = $expertId;
    }

    /**
     * @return int|null
     */
    public function getNotificationType(): ?int
    {
        return $this->notificationType;
    }

    /**
     * @param int|null $notificationType
     */
    public function setNotificationType(?int $notificationType): void
    {
        $this->notificationType = $notificationType;
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
     */
    public function setComplaintId(?int $complaintId): void
    {
        $this->complaintId = $complaintId;
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
     */
    public function setComplaintNumber(?string $complaintNumber): void
    {
        $this->complaintNumber = $complaintNumber;
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
     */
    public function setConsumptionId(?int $consumptionId): void
    {
        $this->consumptionId = $consumptionId;
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
     */
    public function setConsumptionNumber(?string $consumptionNumber): void
    {
        $this->consumptionNumber = $consumptionNumber;
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
     */
    public function setAgreementNumber(?string $agreementNumber): void
    {
        $this->agreementNumber = $agreementNumber;
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
     */
    public function setDate(?string $date): void
    {
        $this->date = $date;
    }

    /**
     * @return string|null
     */
    public function getDifferences(): ?string
    {
        return $this->differences;
    }

    /**
     * @param string|null $differences
     */
    public function setDifferences(?string $differences): void
    {
        $this->differences = $differences;
    }
}
