<?php

namespace NW\WebService\References\Operations\Notification\Validation\Rule;

use NW\WebService\References\Operations\Notification\Contractor;
use NW\WebService\References\Operations\Notification\Validation\ValidatorInterface;

/**
 * ClientValidator class
 */
class ClientValidator implements ValidatorInterface
{

    public function validate(array $data, array &$result = []): bool
    {
        $resellerId = (int)($data['resellerId'] ?? 0);
        $clientId = (int)($data['clientId'] ?? 0);

        $client = Contractor::getById($clientId);

        if (empty($client) || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Client not found!', 400);
        }

        return true;
    }
}
