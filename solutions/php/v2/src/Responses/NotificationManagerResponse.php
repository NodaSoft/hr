<?php

namespace src\Responses;

use Israil\V2\Interfaces\ResponseAbstract;

class NotificationManagerResponse extends ResponseAbstract
{
	// NOTIFICATION MANAGER RESPONSE PROPERTIES

	public bool $success = true;
	public ?string $message = null;

	public function __constructor(array $rawResponse)
	{
		$this->success = $rawResponse['code'] === 200;
		$this->message = $rawResponse['message'];
	}
}