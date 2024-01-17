<?php

namespace App\v2\Notifications;

use App\v2\Events\NotificationEvent;
use App\v2\Facades\BaseNotification;

class NotificationEmployee extends BaseNotification
{
    public function __construct
    (
        NotificationEvent $event
    )
    {
        parent::__construct($event);
    }

    public function send() : bool {
        try {
            $this->event->dispatch();
        } catch (\Exception $ex) {
            $this->error = $ex->getMessage();
            return false;
        }
        return true;
    }

}