<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Exceptions;

use Exception;

/**
 * Class EntityNotExistException
 * Exception thrown while Entity is not found
 */
class EntityNotExistException extends Exception
{
}