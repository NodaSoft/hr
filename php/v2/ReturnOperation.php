<?php


namespace App\v2;

use App\v2\Events\NotificationEvent;
use App\v2\Mail\MessagesClient;
use App\v2\Notifications\NotificationClientEmail;
use App\v2\Notifications\NotificationClientSms;
use App\v2\Notifications\NotificationEmployee;
use App\v2\Responses\NotificationResponse;
use App\v2\Facades\DB;
use App\v2\Requests\NotificationRequest;
use Exception;


class ReturnOperation
{

    /**
     * @return array
     * @throws Exception
     */
    public function doOperation(): array
    {
        $data = (new NotificationRequest($_REQUEST))->validated();
        $response = new NotificationResponse();

        $emailFrom = DB::getResellerEmailFrom();
        $emails = DB::getEmailsByPermit($data['resellerId']);

        $client = DB::getClientById($data['clientId']) ?? throw new Exception('Client not found', 422);

        if ($emailFrom) {
            if (count($emails) > 0) {
                foreach ($emails as $email) {
                    $messageToEmployee = new MessagesClient($data);
                    $response->setNotifyEmployee(new NotificationEmployee(new NotificationEvent($messageToEmployee)));
                }
            }
            if ($client->email) {
                $messageToClientEmail = new MessagesClient($data);
                $response->setClientEmailNotify(new NotificationClientEmail(new NotificationEvent($messageToClientEmail)));
            }
            if ($client->mobile) {
                $messageToClientBySMS = new MessagesClient($data);
                $response->setClientSMSNotify(new NotificationClientSms(new NotificationEvent($messageToClientBySMS)));
            }
        }

        return $response->send();

    }

}
