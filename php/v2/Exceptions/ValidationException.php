<?php

namespace NW\WebService\References\Operations\Notification\Exceptions;

use Exception;

class ValidationException extends Exception
{
  public function __construct(string $message, int $code = 400)
  {
    parent::__construct($message, $code);
  }
}