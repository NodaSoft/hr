<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Struct;

use Symfony\Component\Validator\Constraints as Assert;

final class Template
{
    #[Assert\Positive]
    public readonly int $complaintId;

    #[Assert\NotBlank]
    public readonly string $complaintNumber;

    #[Assert\Positive]
    public readonly int $creatorId;

    #[Assert\NotBlank]
    public readonly string $creatorName;

    #[Assert\Positive]
    public readonly int $expertId;

    #[Assert\NotBlank]
    public readonly string $expertName;

    #[Assert\Positive]
    public readonly int $clientId;

    #[Assert\NotBlank]
    public readonly string $clientName;

    #[Assert\Positive]
    public readonly int $consumptionId;

    #[Assert\NotBlank]
    public readonly string $consumptionNumber;

    #[Assert\NotBlank]
    public readonly string $agreementNumber;

    #[Assert\NotBlank]
    public readonly string $date;


    public readonly ?Differences $differences;

    /**
     * @param positive-int $complaintId
     * @param non-empty-string $complaintNumber
     * @param positive-int $creatorId
     * @param non-empty-string $creatorName
     * @param positive-int $expertId
     * @param non-empty-string $expertName
     * @param positive-int $clientId
     * @param non-empty-string $clientName
     * @param positive-int $consumptionId
     * @param non-empty-string $consumptionNumber
     * @param non-empty-string $agreementNumber
     * @param non-empty-string $date
     */
    public function __construct(
        int $complaintId,
        string $complaintNumber,
        int $creatorId,
        string $creatorName,
        int $expertId,
        string $expertName,
        int $clientId,
        string $clientName,
        int $consumptionId,
        string $consumptionNumber,
        string $agreementNumber,
        string $date,
        ?Differences $differences = null
    ) {
        $this->differences       = $differences;
        $this->date              = $date;
        $this->agreementNumber   = $agreementNumber;
        $this->consumptionNumber = $consumptionNumber;
        $this->consumptionId     = $consumptionId;
        $this->clientName        = $clientName;
        $this->clientId          = $clientId;
        $this->expertName        = $expertName;
        $this->expertId          = $expertId;
        $this->creatorName       = $creatorName;
        $this->creatorId         = $creatorId;
        $this->complaintNumber   = $complaintNumber;
        $this->complaintId       = $complaintId;
    }
}
