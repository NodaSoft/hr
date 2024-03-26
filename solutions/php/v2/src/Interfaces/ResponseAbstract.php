<?php

namespace Israil\V2\Interfaces;

abstract class ResponseAbstract
{
	public function toArray(): array
	{
		return get_object_vars($this);
	}
}