<?php

namespace App\v2\Notifications;

use App\v2\Facades\BaseNotification;
use App\v2\Facades\Event;

class NotificationClientSms extends BaseNotification
{
    public function __construct(Event $event)
    {
        parent::__construct($event);
    }

    public function send(): bool
    {
        try {
            $this->event->dispatch();
        } catch (\Exception $ex) {
            $this->error = $ex->getMessage();
            return false;
        }
        return true;
    }
    public function getMessage(): string {
        return $this->error;
    }
}