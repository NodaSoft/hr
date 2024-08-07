<?php

namespace NW\WebService\References\Operations\Notification\Exceptions;

use Exception;

class ValidationException extends Exception
{
    private const VALIDATION_ERROR_MESSAGE = 'Unprocessable Entity';
    public array $fields;

    /**
     * ValidationException constructor.
     * @param array $fields
     */
    public function __construct(array $fields)
    {
        $this->fields = $fields;
        parent::__construct(self::VALIDATION_ERROR_MESSAGE, 0, null);
    }
}