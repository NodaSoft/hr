<?php

namespace Israil\V2\Responses;

use Israil\V2\Interfaces\ResponseAbstract;

class OperationResponse extends ResponseAbstract
{
	protected bool $notificationEmployeeByEmail = false;
	protected bool $notificationClientByEmail = false;
	protected array $notificationClientBySms = [
		'isSent' => false,
		'message' => '',
	];

	public function setIsEmployeeByEmail(): static
	{
		$this->notificationEmployeeByEmail = true;

		return $this;
	}

	public function setIsClientByEmail(): static
	{
		$this->notificationClientByEmail = true;

		return $this;
	}

	public function setClientBySmsMessage(string $message): static
	{
		$this->notificationClientBySms['message'] = $message;

		return $this;
	}

	public function setClientBySmsIsSent(): static
	{
		$this->notificationClientBySms['isSent'] = true;

		return $this;
	}

	public function errorResponse(string $message): static
	{
		return $this
			->setErrorResponse()
			->setClientBySmsMessage($message);
	}

	private function setErrorResponse(): static
	{
		$this->notificationClientBySms['isSent'] = false;
		$this->notificationClientByEmail = false;
		$this->notificationEmployeeByEmail = false;
		// or something else...

		return $this;
	}
}