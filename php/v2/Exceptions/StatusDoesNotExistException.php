<?php

namespace NW\WebService\References\Operations\Notification\Exceptions;

class StatusDoesNotExistException extends \Exception {
    public function __construct()
    {
        parent::__construct("Status not found", 400);
    }
}
