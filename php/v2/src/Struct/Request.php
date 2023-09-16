<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Struct;

use Exception;
use NW\WebService\References\Operations\Notification\Enum\Status;

final class Request
{
    /**
     * @var positive-int
     */
    public readonly int $resellerId;

    /**
     * @var positive-int
     */
    public readonly int $notificationType;

    /**
     * @var positive-int
     */
    public readonly int $clientId;

    public readonly ?Differences $differences;

    /**
     * @var positive-int
     */
    public readonly int $creatorId;

    /**
     * @var positive-int
     */
    public readonly int $expertId;

    /**
     * @var positive-int
     */
    public readonly int $complaintId;

    /**
     * @var non-empty-string
     */
    public readonly string $complaintNumber;

    /**
     * @var positive-int
     */
    public readonly int $consumptionId;

    /**
     * @var non-empty-string
     */
    public readonly string $consumptionNumber;

    /**
     * @var non-empty-string
     */
    public readonly string $agreementNumber;

    /**
     * @var non-empty-string
     */
    public readonly string $date;

    /**
     * @param positive-int $resellerId
     * @param positive-int $notificationType
     * @param positive-int $clientId
     * @param positive-int $creatorId
     * @param positive-int $expertId
     * @param positive-int $complaintId
     * @param non-empty-string $complaintNumber
     * @param positive-int $consumptionId
     * @param non-empty-string $consumptionNumber
     * @param non-empty-string $agreementNumber
     * @param non-empty-string $date
     * @param non-empty-string[]|null $differences
     * @throws Exception
     */
    public function __construct(
        int $resellerId,
        int $notificationType,
        int $clientId,
        int $creatorId,
        int $expertId,
        int $complaintId,
        string $complaintNumber,
        int $consumptionId,
        string $consumptionNumber,
        string $agreementNumber,
        string $date,
        ?array $differences = null
    ) {
        $this->resellerId = $resellerId;

        $this->notificationType = $notificationType;

        $this->clientId = $clientId;

        $this->creatorId = $creatorId;

        $this->expertId = $expertId;

        $this->complaintId = $complaintId;

        $this->complaintNumber = $complaintNumber;

        $this->consumptionId = $consumptionId;

        $this->consumptionNumber = $consumptionNumber;

        $this->agreementNumber = $agreementNumber;

        $this->date = $date;

        if ($differences !== null && isset($differences['from']) && isset($differences['to'])) {
            $this->differences = new Differences(Status::from($differences['from']), Status::from($differences['to']));
        } else {
            $this->differences = null;
        }
    }
}
