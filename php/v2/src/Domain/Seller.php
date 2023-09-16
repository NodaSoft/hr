<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Domain;

final class Seller extends Contractor
{
    public string $email = 'contractor@example.com';

    /**
     * @return non-empty-string[]
     */
    public function getEmails(): array
    {
        // fakes the method
        return ['someemeil@example.com', 'someemeil2@example.com'];
    }
}
