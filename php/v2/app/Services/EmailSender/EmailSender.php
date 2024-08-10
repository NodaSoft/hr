<?php

namespace app\Services\EmailSender;

use app\Domain\Notification\Exceptions\ActionException;
use app\Services\EmailSender\DTO\EmailDTO;

class EmailSender
{
    public function sendEmail(EmailDTO $data): bool
    {
        $isSend = false;

        try {
            if ($data->client_id) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $data->email_from,
                        'emailTo' => $data->client_email,
                        'subject' => __('complaintClientEmailSubject', $data->data, $data->user_id),
                        'message' => __('complaintClientEmailBody', $data->data, $data->user_id),
                    ],
                ], $data->user_id, $data->client_id, $data->status, $data->to);
            } else {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $data->email_from,
                        'emailTo' => $data->email_to,
                        'subject' => __('complaintEmployeeEmailSubject', $data->data, $data->user_id),
                        'message' => __('complaintEmployeeEmailBody', $data->data, $data->user_id),
                    ],
                ], $data->user_id, $data->status);
            }

            $isSend = true;
        } catch (\Exception $e) {
            $errorMessage = $e->getMessage() . "\n" . $e->getFile() . ":" . $e->getLine() . "\n" . $e->getTraceAsString();
            throw new ActionException($errorMessage);
        }

        return $isSend;
    }
}