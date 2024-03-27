<?php

namespace Israil\V2\Templates;

class EmailTemplate
{
	public ?int $COMPLAINT_ID = null;
	public ?int $CLIENT_ID = null;
	public ?int $CREATOR_ID = null;
	public ?int $EXPERT_ID = null;
	public ?int $CONSUMPTION_ID = null;
	public ?string $COMPLAINT_NUMBER = null;
	public ?string $CREATOR_NAME = null;
	public ?string $EXPERT_NAME = null;
	public ?string $CLIENT_NAME = null;
	public ?string $CONSUMPTION_NUMBER = null;
	public ?string $AGREEMENT_NUMBER = null;
	public array $DIFFERENCES = [];
	public ?string $DATE = null;

	public function toArray(): array
	{
		return get_object_vars($this);
	}
}