<?php

namespace NW\WebService\References\Operations\Notification;

class NotificationManager
{
    /**
     * Send notification
     *
     * @param int $resellerId
     * @param int $clientId
     * @param string $event
     * @param string $newPositionStatus
     * @param array $templateData
     * @param string $error
     * @return bool
     */
    public static function send(int $resellerId, int $clientId, string $event, string $newPositionStatus, array $templateData, string &$error): bool
    {
        if ($result = (bool)mt_rand(0, 1)) {
            $error = '';
        } else {
            $error = 'fake error message';
        }

        return $result;
    }
}