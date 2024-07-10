<?php

namespace NW\WebService\References\Operations\Notification\Exceptions;

class CreatorNotExistException extends \Exception {
    public function __construct()
    {
        parent::__construct("Creator not found", 400);
    }
}
