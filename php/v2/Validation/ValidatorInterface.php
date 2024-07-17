<?php

namespace NW\WebService\References\Operations\Notification\Validation;

use NW\WebService\References\Operations\Notification\Dto\NotificationData;
use NW\WebService\References\Operations\Notification\Notification\Exceptions\ValidationException;

/**
 * Interface ValidatorInterface
 *
 * Defines the contract for validators.
 */
interface ValidatorInterface {
    /**
     * @throws ValidationException
     */
    public function validate(NotificationData $data): void;
}