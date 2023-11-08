<?php

namespace NodaSoft\Dependencies;

use NodaSoft\Mail\Mail;

class Dependencies
{
    public function getMail(): Mail
    {
        return new Mail();
    }
}
