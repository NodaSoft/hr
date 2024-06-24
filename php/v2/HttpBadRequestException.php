<?php

namespace NW\WebService\References\Operations\Notification;

use Throwable;

class HttpBadRequestException extends \Exception
{
    public function __construct($message = '', ?Throwable $previous = null)
    {
        parent::__construct($message, Response::HTTP_BAD_REQUEST, $previous);
    }
}