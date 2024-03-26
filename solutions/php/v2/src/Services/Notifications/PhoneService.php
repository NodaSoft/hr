<?php

namespace Israil\V2\Services\Notifications;

use Israil\V2\Stubs\NotificationManager;
use Israil\V2\Templates\EmailTemplate;
use NW\WebService\References\Operations\Notification\NotificationEvents;
use src\Responses\NotificationManagerResponse;

class PhoneService
{
	public function sendSmsClient(
		int           $resellerId,
		int           $clientId,
		string        $type, // to enum
		EmailTemplate $validatedTemplateData,
	): NotificationManagerResponse {
		$rawResponse = NotificationManager::send(
			$resellerId,
			$clientId,
			NotificationEvents::CHANGE_RETURN_STATUS,
			$validatedTemplateData->DIFFERENCES['to'],
			$validatedTemplateData,
			$error,
		);

		return new NotificationManagerResponse($rawResponse);
	}
}