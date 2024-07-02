<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Exception;

use RuntimeException;
use Throwable;

/**
 * Class SelerNotFoundException.
 *
 */
final class SelerNotFoundException extends RuntimeException
{
    /**
     * @inerhitDoc
     */
    public function __construct(string $message = 'Seller not found!', int $code = 400, Throwable $previous = null)
    {
        parent::__construct($message, $code, $previous);
    }
}