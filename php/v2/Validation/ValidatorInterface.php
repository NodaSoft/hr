<?php

namespace NW\WebService\References\Operations\Notification\Validation;

/**
 * Interface ValidatorInterface
 */
interface ValidatorInterface
{

    /**
     * Validates data and writes errors to $result.
     *
     * @param array $data
     * @param array $result
     * @return bool
     * @throws \Exception
     */
    public function validate(array $data, array &$result = []): bool;
}
