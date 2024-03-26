<?php

namespace Israil\V2\Exceptions;

use Exception;

class NotFoundException extends Exception
{
	public function __construct(
		protected $message = "",
	) {
		parent::__construct($this->message, 400);
	}
}