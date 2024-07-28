<?php

namespace NW\WebService\References\Operations\Notification\Validation\Rule;

use NW\WebService\References\Operations\Notification\Seller;
use NW\WebService\References\Operations\Notification\Validation\ValidatorInterface;

/**
 * SellerValidator class
 */
class SellerValidator implements ValidatorInterface
{

    public function validate(array $data, array &$result = []): bool
    {
        $resellerId = (int)($data['resellerId'] ?? 0);
        $reseller = Seller::getById($resellerId);

        if (empty($reseller)) {
            throw new \Exception('Seller not found!', 400);
        }

        return true;
    }
}
