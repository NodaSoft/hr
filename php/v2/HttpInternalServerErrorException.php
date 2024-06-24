<?php

namespace NW\WebService\References\Operations\Notification;

use Throwable;

class HttpInternalServerErrorException extends \Exception
{
    public function __construct($message = '', ?Throwable $previous = null)
    {
        parent::__construct($message, Response::HTTP_INTERNAL_SERVER_ERROR, $previous);
    }
}