<?php

namespace app\Services\SmsSender;

use app\Domain\Notification\Exceptions\ActionException;
use app\Services\SmsSender\DTO\SmsDTO;

class SmsSender
{
    public function send(SmsDTO $data): array
    {
        $isSend = [
            'isSent' => false,
            'error' => '',
        ];

        try {
            NotificationManager::send($data->user_id, $data->client_id, $data->status, $data->to, $data->data);
            $isSend['isSent'] = true;
        } catch (\Exception $e) {
            $errorMessage = $e->getMessage() . "\n" . $e->getFile() . ":" . $e->getLine() . "\n" . $e->getTraceAsString();
            $isSend['error'] = $errorMessage;
        }

        return $isSend;
    }
}