<?php

namespace App\v2\Events;

use App\v2\Facades\Event;
use App\v2\Facades\MailSender;

class NotificationEvent extends Event
{
    public const string CHANGE_RETURN_STATUS = 'changeReturnStatus';
    public const string NEW_RETURN_STATUS = 'newReturnStatus';

    private MailSender $listener;

    public function __construct(MailSender $eventListener)
    {
        $this->listener = $eventListener;
    }

    #[\Override] public function dispatch(...$args): void
    {
        // TODO: Implement dispatch() method.
        $this->listener::sendMessage(...$args);
    }
}
