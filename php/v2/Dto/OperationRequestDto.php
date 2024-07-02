<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Dto;

/**
 * Class OperationRequestDto.
 *
 */
class OperationRequestDto
{
    /**
     * @var int
     */
    private int $resellerId;
    /**
     * @var int
     */
    private int $notificationType;
    /**
     * @var int
     */
    private int $clientId;
    /**
     * @var int
     */
    private int $creatorId;
    /**
     * @var int
     */
    private int $expertId;
    /**
     * @var array|null
     */
    private ?array $differences;
    /**
     * @var int
     */
    private int $complaintId;
    /**
     * @var string
     */
    private string $complaintNumber;
    /**
     * @var int
     */
    private int $consumptionId;
    /**
     * @var string
     */
    private string $consumptionNumber;
    /**
     * @var string
     */
    private string $agreementNumber;
    /**
     * @var string
     */
    private string $date;

    /**
     * OperationRequestDto constructor.
     *
     * @param array $data
     */
    public function __construct(
        array $data = []
    )
    {
        $this->filDto($data);
    }

    /**
     *  fill OperationRequestDto object from array.
     *
     * @param array $data
     */
    public function filDto(array $data)
    {
        foreach ($data as $key => $value) {
            if (property_exists($this, $key)) {
                $this->{$key} = $value;
            }
        }
    }

    /**
     * get ResellerId.
     *
     * @return int|null
     */
    public function getResellerId(): ?int
    {
        return $this->resellerId ?? null;
    }

    /**
     * get NotificationType.
     *
     * @return int|null
     */
    public function getNotificationType(): ?int
    {
        return $this->notificationType ?? null;
    }

    /**
     * get ClientId.
     *
     * @return int|null
     */
    public function getClientId(): ?int
    {
        return $this->clientId ?? null;
    }

    /**
     * get creatorId.
     *
     * @return int|null
     */
    public function getCreatorId(): ?int
    {
        return $this->creatorId ?? null;
    }

    /**
     * get expertId.
     *
     * @return int|null
     */
    public function getExpertId(): ?int
    {
        return $this->expertId ?? null;
    }

    /**
     * get differences.
     *
     * @return array
     */
    public function getDifferences(): array
    {
        return $this->differences ?? [];
    }

    /**
     * get complaintId.
     *
     * @return int|null
     */
    public function getComplaintId(): ?int
    {
        return $this->complaintId ?? null;
    }

    /**
     * get complaintNumber.
     *
     * @return string|null
     */
    public function getComplaintNumber(): ?string
    {
        return $this->complaintNumber ?? null;
    }

    /**
     * get complaintNumber.
     *
     * @return int|null
     */
    public function getConsumptionId(): ?int
    {
        return $this->consumptionId ?? null;
    }

    /**
     * get complaintNumber.
     *
     * @return string|null
     */
    public function getConsumptionNumber(): ?string
    {
        return $this->consumptionNumber ?? null;
    }

    /**
     * get complaintNumber.
     *
     * @return string|null
     */
    public function getAgreementNumber(): ?string
    {
        return $this->agreementNumber ?? null;
    }

    /**
     * Convert OperationRequestDto to array.
     *
     * @return string|null
     */
    public function getDate(): ?string
    {
        return $this->date ?? null;
    }

}