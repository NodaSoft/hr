<?php

namespace NW\WebService\References\Operations\Notification\Mailer;

/**
 * MailableTrait trait
 */
trait MailableTrait
{
    /**
     * Returns the reseller's email address.
     *
     * @return string
     */
    public function getResellerEmailFrom(): string
    {
        return 'contractor@example.com';
    }

    /**
     * Returns a list of emails permitted for a specific event.
     *
     * @param int $resellerId
     * @param string $event
     * @return array
     */
    public function getEmailsByPermit(int $resellerId, string $event): array
    {
        // fakes the method
        return ['someemail@example.com', 'someemail2@example.com'];
    }
}
