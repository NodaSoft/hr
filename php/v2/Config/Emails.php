<?php

declare(strict_types=1);


namespace NW\WebService\Config;

class Emails
{
    public static function getEmailsByPermit(): array
    {
        return ['someemeil@example.com', 'someemeil2@example.com'];
    }
}