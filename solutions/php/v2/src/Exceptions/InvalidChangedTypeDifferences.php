<?php

namespace Israil\V2\Exceptions;

class InvalidChangedTypeDifferences extends \Exception
{
	public function __construct()
	{
		parent::__construct("Invalid or NULL: DIFFERENCES[FROM AND\OR]", 500);
	}
}