<?php

namespace App\v2\Responses;

use App\v2\Notifications\NotificationClientEmail;
use App\v2\Notifications\NotificationClientSms;
use App\v2\Notifications\NotificationEmployee;

class NotificationResponse
{
    public function __construct
    (
        private ?NotificationEmployee    $employeeNotify = null,
        private ?NotificationClientEmail $clientNotify = null,
        private ?NotificationClientSms   $clientSmsNotify = null,
    )
    {
    }

    public function setNotifyEmployee(NotificationEmployee $employee) : void {
        $this->employeeNotify = $employee;
    }

    public function setClientEmailNotify(NotificationClientEmail $clientEmail) : void {
        $this->clientNotify = $clientEmail;
    }
    public function setClientSMSNotify(NotificationClientSms $clientSms) : void {
        $this->clientSmsNotify = $clientSms;
    }

    public function send
    (
    ) : array {
        return [
            'notificationEmployeeByEmail' => $this->employeeNotify && $this->employeeNotify->send
                (
                ),
            'notificationClientByEmail' => $this->clientNotify && $this->clientNotify->send
                (
                ),
            'notificationClientBySms' => [
                'isSent' => $this->clientSmsNotify && $this->clientSmsNotify->send
                    (
                    ),
                'message' => $this->clientSmsNotify ? $this->clientSmsNotify->getMessage() : '',
                ]
        ];
    }
}