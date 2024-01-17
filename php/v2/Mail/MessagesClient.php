<?php

namespace App\v2\Mail;

use App\v2\Exceptions\TemplateKeyException;
use App\v2\Facades\MailSender;


class MessagesClient extends MailSender
{
    private array $templateData = [];

    /**
     * @throws \Exception
     */
    public function __construct(array $data)
    {
        foreach ($data as $key => $value) {
            $this->templateData[strtoupper($key)] = $value;
            if (!$value) {
                throw new TemplateKeyException("Template Data ({$key}) is empty!", 500);
            }
        }
    }

    /**
     * @param string $emailFrom
     * @param string $emailTo
     * @param int $resellerId
     * @param int $statusId
     * @param int $clientId
     * @return void
     */
    public static function sendMessage(
        string $emailFrom,
        string $emailTo,
        int $resellerId,
        int $statusId,
        int $clientId
    ): void
    {
        /** TODO: email sender */
    }

}
