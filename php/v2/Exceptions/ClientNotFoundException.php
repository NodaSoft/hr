<?php

namespace NW\WebService\References\Operations\Notification\Exceptions;

class ClientNotFoundException extends \Exception {
    public function __construct()
    {
        parent::__construct("Client not found", 400);
    }
}
