<?php

namespace Israil\V2\Interfaces;

interface ValidatorInterface
{
	public static function validate(object $object): object;
}