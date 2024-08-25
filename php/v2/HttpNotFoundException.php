<?php

namespace NW\WebService\References\Operations\Notification;

use Throwable;

class HttpNotFoundException extends \Exception
{
    public function __construct($message = '', ?Throwable $previous = null)
    {
        parent::__construct($message, Response::HTTP_NOT_FOUND, $previous);
    }
}