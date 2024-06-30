<?php

namespace NW\WebService\References\Operations\Notification\Dto;

readonly class ReturnOperationRequestDto
{
    /**
     * @param int $resellerId
     * @param int $notificationType
     * @param int $complaintId
     * @param int $complaintNumber
     * @param int $creatorId
     * @param int $expertId
     * @param int $clientId
     * @param int $consumptionId
     * @param string $consumptionNumber
     * @param string $agreementNumber
     * @param string $date
     * @param array{from: int, to: int}|null $differences
     */
    public function __construct(
        public int $resellerId,
        public int $notificationType,
        public int $complaintId,
        public int $complaintNumber,
        public int $creatorId,
        public int $expertId,
        public int $clientId,
        public int $consumptionId,
        public string $consumptionNumber,
        public string $agreementNumber,
        public string $date,
        public ?array $differences,
    ) {
    }
}
