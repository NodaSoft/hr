<?php

namespace NW\WebService\References\Operations\Notification\Notification\Exceptions;

use Exception;
use Throwable;

/**
 * Пользовательское исключение для ошибок валидации
 */
class ValidationException extends Exception
{
    public function __construct(string $message = "", int $code = 0, ?Throwable $previous = null)
    {
        parent::__construct($message, $code, $previous);
    }
}