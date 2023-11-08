<?php

namespace NodaSoft\Mail;

class Mail
{
    public function send(Message $message): Result
    {
        $result = new Result($message->getRecipient());
        try {
            $isSent = $this->mail(
                $message->getTo(),
                $message->getSubject(),
                $message->getMessage(),
                $message->getHeaders(),
                $message->getParams()
            );
            $result->setIsSent($isSent);
        } catch (\Throwable $th) {
            $result->setErrorMessage($th->getMessage());
        }

        return $result;
    }

    public function mail(
        string $to,
        string $subject,
        string $message,
        string $headers,
        string $params
    ): bool
    {
        return mail(
            $to(),
            $subject(),
            $message(),
            $headers(),
            $params()
        );
    }
}
