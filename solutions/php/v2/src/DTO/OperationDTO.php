<?php

namespace Israil\V2\Dto;

use Israil\V2\Enums\NotificationTypeEnum;
use Israil\V2\Interfaces\DtoInterface;

class OperationDTO implements DtoInterface
{
	public function __construct(
		public ?string               $resellerId,
		public NotificationTypeEnum $notificationType,
		public ?int                  $clientId,
		public ?int                  $creatorId,
		public ?int                  $expertId,
		public ?string               $complaintId,
		public ?string               $complaintNumber,
		public ?int                  $consumptionId,
		public ?string               $consumptionNumber,
		public ?string               $agreementNumber,
		public ?array                $differences,
		public ?string               $date,
	)
	{
	}

	// getters setters...

	public function isType(NotificationTypeEnum $type): bool
	{
		return $this->notificationType === $type;
	}

	public static function fromArray(array $data): static
	{
		return new static(
			$data['resellerId'] ?? null,
			NotificationTypeEnum::from($data['notificationType']),
			$data['clientId'] ?? null,
			$data['creatorId'] ?? null,
			$data['expertId'] ?? null,
			$data['complaintId'] ?? null,
			$data['complaintNumber'] ?? null,
			$data['consumptionId'] ?? null,
			$data['consumptionNumber'] ?? null,
			$data['agreementNumber'] ?? null,
			$data['differences'] ?? null,
			$data['date'] ?? null,
		);
	}
}