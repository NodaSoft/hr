<?php

namespace NW\WebService\References\Operations\Notification\Exceptions;

class EmptyNotificationTypeException extends \Exception {
    public function __construct()
    {
        parent::__construct('Empty notificationType', 400);
    }
}
