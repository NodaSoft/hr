<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Contracts;


/**
 * Interface NotificationManagerInterface.
 *
 */
interface NotificationManagerInterface
{
    /**
     * @param int $resellerId
     * @param int $clientId
     * @param string $notificationType
     * @param int|null $differencesTo
     * @param array $templateData
     * @return void
     */
    public function send(
        int    $resellerId,
        int    $clientId,
        string $notificationType,
        int    $differencesTo,
        array  $templateData
    ): bool;
}