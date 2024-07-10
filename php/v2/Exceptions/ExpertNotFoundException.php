<?php

namespace NW\WebService\References\Operations\Notification\Exceptions;

class ExpertNotFoundException extends \Exception {
    public function __construct()
    {
        parent::__construct("Expert not found", 400);
    }
}
