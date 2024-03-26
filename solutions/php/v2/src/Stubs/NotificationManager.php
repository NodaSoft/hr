<?php

namespace Israil\V2\Stubs;

use Israil\V2\Templates\EmailTemplate;

class NotificationManager
{
	public static function send(
		int $resellerId,
		int $clientId,
		string $type,
		string $differencesTo,
		EmailTemplate $template,
		string &$error,
	): array {
		return ['code' => 200, 'message' => 'Success'];
	}
}