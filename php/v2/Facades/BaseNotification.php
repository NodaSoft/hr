<?php

namespace App\v2\Facades;


abstract class BaseNotification
{
    protected string $error = '';
    public function __construct
    (
        protected Event $event
    )
    {
    }

}