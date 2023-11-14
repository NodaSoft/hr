<?php

namespace NodaSoft\ReferencesOperation\Params;

use NodaSoft\Request\Request;

class ReturnOperationStatusChangedParams implements ReferencesOperationParams
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

    /** @var int */
    private $previousStatusId;

    /** @var int */
    private $currentStatusId;

    public function setRequest(Request $request): void
    {
        foreach ($this as $key => $value) {
            $setter = 'set' . $key;
            if (method_exists($this, $setter)) {
                $this->$setter($request->getData($key));
            }
            $this->setDifferences($request->getData('differences'));
        }
    }

    public function isValid(): bool
    {
        if (empty($this->resellerId)) {
            return false;
        }

        if (empty($this->notificationType)) {
            return false;
        }

        if (empty($this->creatorId)) {
            return false;
        }

        if (empty($this->expertId)) {
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

    public function getResellerId(): ?int
    {
        return $this->resellerId;
    }

    public function setResellerId(?int $resellerId): void
    {
        $this->resellerId = $resellerId;
    }

    public function getClientId(): ?int
    {
        return $this->clientId;
    }

    public function setClientId(?int $clientId): void
    {
        $this->clientId = $clientId;
    }

    public function getCreatorId(): ?int
    {
        return $this->creatorId;
    }

    public function setCreatorId(?int $creatorId): void
    {
        $this->creatorId = $creatorId;
    }

    public function getExpertId(): ?int
    {
        return $this->expertId;
    }

    public function setExpertId(?int $expertId): void
    {
        $this->expertId = $expertId;
    }

    public function getNotificationType(): ?int
    {
        return $this->notificationType;
    }

    public function setNotificationType(?int $notificationType): void
    {
        $this->notificationType = $notificationType;
    }

    public function getComplaintId(): ?int
    {
        return $this->complaintId;
    }

    public function setComplaintId(?int $complaintId): void
    {
        $this->complaintId = $complaintId;
    }

    public function getComplaintNumber(): ?string
    {
        return $this->complaintNumber;
    }

    public function setComplaintNumber(?string $complaintNumber): void
    {
        $this->complaintNumber = $complaintNumber;
    }

    public function getConsumptionId(): ?int
    {
        return $this->consumptionId;
    }

    public function setConsumptionId(?int $consumptionId): void
    {
        $this->consumptionId = $consumptionId;
    }

    public function getConsumptionNumber(): ?string
    {
        return $this->consumptionNumber;
    }

    public function setConsumptionNumber(?string $consumptionNumber): void
    {
        $this->consumptionNumber = $consumptionNumber;
    }

    public function getAgreementNumber(): ?string
    {
        return $this->agreementNumber;
    }

    public function setAgreementNumber(?string $agreementNumber): void
    {
        $this->agreementNumber = $agreementNumber;
    }

    public function getDate(): ?string
    {
        return $this->date;
    }

    public function setDate(?string $date): void
    {
        $this->date = $date;
    }

    public function getPreviousStatusId(): ?int
    {
        return $this->previousStatusId;
    }

    public function getCurrentStatusId(): ?int
    {
        return $this->currentStatusId;
    }

    public function setDifferences(?array $differences): void
    {
        $this->previousStatusId = $differences['from']
            ? (int) $differences['from']
            : null;

        $this->currentStatusId = $differences['to']
            ? (int) $differences['to']
            : null;
    }
}
