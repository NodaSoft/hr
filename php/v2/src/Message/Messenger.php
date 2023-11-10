<?php

namespace NodaSoft\Message;

class Messenger
{
    /** @var Client */
    private $client;

    public function __construct(Client $client)
    {
        $this->client = $client;
    }

    public function send(Message $message): Result
    {
        $result = new Result($message->getRecipient(), static::class);

        if (! $this->client->isValid($message)) {
            $result->setErrorMessage("Invalid parameters. Is failed to send a message.");
            return $result;
        }

        try {
            $isSent = $this->client->send($message);
            $result->setIsSent($isSent);
        } catch (\Throwable $th) {
            $result->setErrorMessage($th->getMessage());
        }

        return $result;
    }
}
