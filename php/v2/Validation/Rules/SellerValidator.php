<?php

namespace NW\WebService\References\Operations\Notification\Validation;

use NW\WebService\References\Operations\Notification\Dto\NotificationData;
use NW\WebService\References\Operations\Notification\Notification\Exceptions\ValidationException;
use NW\WebService\References\Operations\Notification\Seller;

/**
 * Валидатор для проверки продавца
 */
readonly class SellerValidator implements ValidatorInterface
{
    public function __construct(private Seller $seller)
    {
    }

    public function validate(NotificationData $data): void
    {
        if ($this->seller->getById($data->resellerId) === null) {
            throw new ValidationException('Seller not found!');
        }
    }
}