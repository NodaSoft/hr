<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Contracts;

use NW\WebService\References\Operations\Notification\Contracts\MessagesClientInterface;

/**
 * Class MessagesClient.
 *
 */
class MessagesClient implements MessagesClientInterface
{
    /**
     * @param int $resellerId
     * @return string
     */
    function getResellerEmailFrom(int $resellerId): string
    {
        return 'contractor@example.com';
    }

    /**
     * @param int $resellerId
     * @param string $event
     * @return array
     */
    function getEmailsByPermit(int $resellerId, string $event): array
    {
        // fakes the method
        return ['someemeil@example.com', 'someemeil2@example.com'];
    }

    /**
     * @inerhitDoc
     */
    public function sendMessages(
        array  $messages,
        int    $resellerId,
        ?int   $clientId,
        string $notificationType,
        ?int   $differencesTo = null
    ): void
    {
        // do some stuff
    }
}