<?php

namespace NW\WebService\References\Operations\Notification\Validation;

use NW\WebService\References\Operations\Notification\Contractor;
use NW\WebService\References\Operations\Notification\Dto\NotificationData;
use NW\WebService\References\Operations\Notification\Notification\Exceptions\ValidationException;

/**
 * Валидатор для проверки клиента
 */
readonly class ClientValidator implements ValidatorInterface
{
    public function __construct(private Contractor $contractor) {}

    public function validate(NotificationData $data): void
    {
        $client = $this->contractor->getById($data->clientId);

        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->sellerId !== $data->resellerId) {
            throw new ValidationException('Client not found!');
        }
    }
}
