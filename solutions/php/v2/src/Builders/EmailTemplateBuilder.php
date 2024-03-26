<?php

namespace Israil\V2\Builders;

use Israil\V2\Enums\NotificationTypeEnum;
use Israil\V2\Interfaces\BuilderInterface;
use Israil\V2\Templates\EmailTemplate;
use NW\WebService\References\Operations\Notification\Status;

class EmailTemplateBuilder implements BuilderInterface
{
	public function __construct(
		private EmailTemplate $template
	)
	{
	}

	public function setComplaintId(int $id): static
	{
		$this->template->COMPLAINT_ID = $id;

		return $this;
	}

	public function setComplaintNumber(string $number): static
	{
		$this->template->COMPLAINT_NUMBER = $number;

		return $this;
	}

	public function setClientId(string $id): static
	{
		$this->template->CLIENT_ID = $id;

		return $this;
	}

	public function setCreatorId(string $id): static
	{
		$this->template->CREATOR_ID = $id;

		return $this;
	}

	public function setExpertId(int $id): static
	{
		$this->template->EXPERT_ID = $id;

		return $this;
	}

	public function setConsumptionId(int $id): static
	{
		$this->template->CONSUMPTION_ID = $id;

		return $this;
	}

	public function setCreatorName(string $name): static
	{
		$this->template->CREATOR_NAME = $name;

		return $this;
	}

	public function setExpertName(string $name): static
	{
		$this->template->EXPERT_NAME = $name;

		return $this;
	}

	public function setClientName(string $name): static
	{
		$this->template->CLIENT_NAME = $name;

		return $this;
	}

	public function setConsumptionNumber(string $number): static
	{
		$this->template->CONSUMPTION_NUMBER = $number;

		return $this;
	}

	public function setAgreementNumber(string $number): static
	{
		$this->template->AGREEMENT_NUMBER = $number;

		return $this;
	}

	public function setDifferences(NotificationTypeEnum $notificationType, int $resellerId, ?array $differences): static
	{
		if ($notificationType === NotificationTypeEnum::NEW) {
			$this->template->DIFFERENCES = __('NewPositionAdded', null, $resellerId);
		} elseif ($notificationType === NotificationTypeEnum::CHANGE &&
			!empty($differences['from']) &&
			!empty($differences['to'])
		) {
			$this->template->DIFFERENCES = __('PositionStatusHasChanged', [
				'FROM' => Status::getName((int) $differences['from']),
				'TO' => Status::getName((int) $differences['to']),
			], $resellerId);
		} else {
			$this->template->DIFFERENCES = [];
		}

		return $this;
	}

	public function setDate(string $date): static
	{
		$this->template->DATE = $date;

		return $this;
	}

	public function resolve()
	{
		return $this->template;
	}
}