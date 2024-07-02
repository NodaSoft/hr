<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Helpers;

use NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException;

/**
 * Class NotificationsManagement
 * Manages Contractors Notifications
 */
class NotificationsManagement
{
    /** @var int */
    public const MESSAGE_TYPE_EMAIL = 0;

    /**
     * @param string $from
     * @param array|string $to
     * @param array $templateData
     * @param string $subjectTemplate
     * @param string $messageTemplate
     * @return bool
     * @throws \NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException
     */
    public function emailSend(
        string $from,
        array|string $to,
        array $templateData,
        string $subjectTemplate,
        string $messageTemplate
    ): bool {
        $result = false;
        $this->validateEmailData($templateData);
        $to = array_filter(array_unique(is_array($to) ? $to : [$to]));
        if ($from && count($to)) {
            foreach ($to as $email) {
                MessagesClient::sendMessage(
                    [
                        self::MESSAGE_TYPE_EMAIL => [
                            'emailFrom' => $from,
                            'emailTo' => $email,
                            'subject' => Support::__($subjectTemplate, $templateData),
                            'message' => Support::__($messageTemplate, $templateData),
                        ],
                    ]
                );
            }
            $result = true;
        }

        return $result;
    }

    /**
     * @param string|null $mobile
     * @param int $resellerId
     * @param int $clientId
     * @param string $event
     * @param array $templateData
     * @return bool
     * @throws \NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException
     */
    public function smsSend(?string $mobile, int $resellerId, int $clientId, string $event, array $templateData): bool
    {
        $result = false;
        if ($mobile) {
            $this->validateEmailData($templateData);

            NotificationManager::send(
                $resellerId,
                $clientId,
                $event,
                $mobile,
                $templateData,
            );

            $result = true;
        }

        return $result;
    }


    /**
     * @param array $templateData
     * @return void
     * @throws \NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException
     */
    private function validateEmailData(array $templateData): void
    {
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new InvalidArgumentsException(
                    Support::__('Template Data (:key) is empty!', [':key' => $key]),
                    500
                );
            }
        }
    }
}