<?php

namespace Israil\V2\Validator;

use Israil\V2\Exceptions\EmailTemplateKeysNotDefined;
use Israil\V2\Interfaces\ValidatorInterface;
use Israil\V2\Templates\EmailTemplate;
use Israil\V2\Exceptions\InvalidChangedTypeDifferences;

class EmailTemplateValidator implements ValidatorInterface
{
	/**
	 * @param EmailTemplate $object
	 * @return EmailTemplate
	 * @throws EmailTemplateKeysNotDefined
	 * @throws InvalidChangedTypeDifferences
	 */
	public static function validate(object $object): object
	{
		self::validateEmptyKeys($object);
		self::validateDifferences($object->DIFFERENCES);

		return $object;
	}

	/**
	 * @throws EmailTemplateKeysNotDefined
	 */
	private static function validateEmptyKeys(EmailTemplate $template): void
	{
		$templateData = get_object_vars($template);
		$emptyKeys = [];

		foreach ($templateData as $key => $tempData) {
			if (empty($tempData)) $emptyKeys[] = $key;
		}

		if (!empty($emptyKeys)) {
			throw new EmailTemplateKeysNotDefined(implode(', ', $emptyKeys));
		}
	}

	/**
	 * @throws InvalidChangedTypeDifferences
	 */
	private static function validateDifferences(array $differences): void
	{
		if (empty($differences['from']) || empty($differences['to'])) {
			throw new InvalidChangedTypeDifferences;
		}
	}
}