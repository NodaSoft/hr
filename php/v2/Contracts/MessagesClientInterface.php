<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Contracts;

use NW\WebService\References\Operations\Notification\Dto\EmailMessageDto;

/**
 * Interface MessagesClientInterface.
 *
 */
interface MessagesClientInterface
{
    /**
     * @param EmailMessageDto[] $messages
     * @param int $resellerId
     * @param int|null $clientId
     * @param string $notificationType
     * @param int|null $differencesTo
     * @return void
     */
    public function sendMessages(
        array  $messages,
        int    $resellerId,
        int    $clientId,
        string $notificationType,
        ?int   $differencesTo = null
    ): void;
}