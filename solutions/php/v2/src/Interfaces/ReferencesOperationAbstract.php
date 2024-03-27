<?php

namespace Israil\V2\Interfaces;

// STUB
abstract class ReferencesOperationAbstract
{
	abstract public function doOperation(): array;

	public function getRequest($pName)
	{
		return $_REQUEST[$pName];
	}
}