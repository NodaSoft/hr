<?php

namespace Israil\V2\Validator;

use Israil\V2\Interfaces\ValidatorInterface;
use Israil\V2\Dto\OperationDTO;
use Israil\V2\Exceptions\NotFoundException;

class OperationRequestValidator implements ValidatorInterface
{
	/**
	 * @param OperationDTO $object
	 * @return OperationDTO
	 * @throws NotFoundException
	 */
	public static function validate(object $object): object
	{
		if (is_null($object->resellerId))
			throw new NotFoundException('Empty resellerId');
		if (is_null($object->notificationType))
			throw new NotFoundException('Empty notificationType');

		return $object;
	}
}