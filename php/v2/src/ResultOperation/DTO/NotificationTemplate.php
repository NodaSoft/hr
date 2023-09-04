<?php

declare(strict_types=1);

namespace ResultOperation\DTO;

/**
 * По хорошему нужно задекларировать геттеры и сеттеры на все
 * параметры шаблона (COMPLAINT_NUMBER, CREATOR_ID и т.д.), но очень лень, извините
 */
class NotificationTemplate
{
    /**
     * @var int
     */
    private int $complaintId;

    /**
     * @var int
     */
    private int $resellerId;

    /**
     * @return int
     */
    public function getComplaintId(): int
    {
        return $this->complaintId;
    }

    /**
     * @param int $complaintId
     * @return self
     */
    public function setComplaintId(int $complaintId): self
    {
        $this->complaintId = $complaintId;

        return $this;
    }

    /**
     * @param int $resellerId
     * @return $this
     */
    public function setResellerId(int $resellerId): self
    {
        $this->resellerId = $resellerId;

        return $this;
    }

    /**
     * @return int
     */
    public function getResellerId(): int
    {
        return $this->resellerId;
    }

    /**
     * @return array
     */
    public function toArray(): array
    {
        return [
            'COMPLAINT_ID'       => $this->complaintId,
//            'COMPLAINT_NUMBER'   => $this->complaintNumber,
//            'CREATOR_ID'         => $this->creatorId,
//            'CREATOR_NAME'       => $this->creatorName,
//            'EXPERT_ID'          => $this->expertId,
//            'EXPERT_NAME'        => $this->expertName,
//            'CLIENT_ID'          => $this->clientId,
//            'CLIENT_NAME'        => $this->clientName,
//            'CONSUMPTION_ID'     => $this->consumptionId,
//            'CONSUMPTION_NUMBER' => $this->consumptionNumber,
//            'AGREEMENT_NUMBER'   => $this->agreementNumber,
//            'DATE'               => $this->date,
//            'DIFFERENCES'        => $this->differences,
        ];
    }
}
