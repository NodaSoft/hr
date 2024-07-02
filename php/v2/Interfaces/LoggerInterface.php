<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Interfaces;

/**
 * Interface LoggerInterface
 * Mock Interface instead Psr Logger
 */
interface LoggerInterface
{
    /**
     * @param string $message
     * @param array $additionalData
     * @return void
     */
    public function log(string $message, array $additionalData): void;
}