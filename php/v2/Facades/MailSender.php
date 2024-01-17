<?php

namespace App\v2\Facades;

abstract class MailSender
{
    public abstract static function sendMessage(
        string $emailFrom,
        string $emailTo,
        int $resellerId,
        int $statusId,
        int $clientId,
    ): void;

}