<?php

namespace NodaSoft\Messenger;

use NodaSoft\Messenger\Recipient;

class Messenger
{
    /** @var Client */
    private $client;

    public function __construct(Client $client)
    {
        $this->client = $client;
    }

    public function send(
        Message   $message,
        Recipient $recipient,
        Recipient $sender
    ): Result {
        $result = new Result($recipient, get_class($this->client));

        if (! $this->client->isValid($message, $recipient, $sender)) {
            $result->setErrorMessage("Invalid parameters. Is failed to send a message.");
            return $result;
        }

        try {
            $isSent = $this->client->send($message, $recipient, $sender);
            $result->setIsSent($isSent);
        } catch (\Throwable $th) {
            $result->setErrorMessage($th->getMessage());
        }

        return $result;
    }
}
