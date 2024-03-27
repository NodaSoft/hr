<?php

namespace Israil\V2\Exceptions;

use Exception;

class EmailTemplateKeysNotDefined extends Exception
{
	public function __construct(string $keyName)
	{
		parent::__construct("Template data keys: $keyName is empty!", 500);
	}
}