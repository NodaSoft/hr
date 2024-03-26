<?php

namespace Israil\V2\Stubs;

// STUB
class Contractor
{
	const TYPE_CUSTOMER = 0;
	public $id;
	public $type;
	public $name;
	public $email = 'random-email@email';
	public $mobile;

	public static function getById(int $resellerId): ?self
	{
		return new self($resellerId); // fakes the getById method
	}

	public function getFullName(): string
	{
		return $this->name . ' ' . $this->id;
	}

	public function canReceiveSms(): bool
	{
		return $this->mobile;
	}
}
