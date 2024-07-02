<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 * Handles Mock classes
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Helpers;

use NW\WebService\References\Operations\Notification\DTO\Notification\OperationResultDTO;

class Support
{
    public static function __(string $string, $arguments = []): string
    {
        return strtr($string, $arguments);
    }
}

class Status
{
    /** @var int */
    public int $id;

    /** @var string */
    public string $name;

    /**
     * @param int $id
     * @return string
     */
    public static function getName(int $id): string
    {
        $a = [
            0 => 'Completed',
            1 => 'Pending',
            2 => 'Rejected',
        ];

        return $a[$id];
    }
}

abstract class ReferencesOperation
{
    /**
     * @return \NW\WebService\References\Operations\Notification\DTO\Notification\OperationResultDTO
     */
    abstract public function doOperation(): OperationResultDTO;

    /**
     * @param string $pName
     * @return mixed
     */
    public function getRequest(string $pName): mixed
    {
        return $_REQUEST[$pName];
    }
}

class NotificationEvents
{
    /** @var string */
    const CHANGE_RETURN_STATUS = 'changeReturnStatus';
    const NEW_RETURN_STATUS    = 'newReturnStatus';
}


class MessagesClient
{
    /**
     * @param array $data
     * @return void
     */
    public static function sendMessage(array $data): void
    {
        //do nothing
    }
}

class NotificationManager
{
    /**
     * @param int $resellerId
     * @param int $clientId
     * @param string $event
     * @param string $mobile
     * @param array $templateData
     * @return void
     */
    public static function send(
        int $resellerId,
        int $clientId,
        string $event,
        string $mobile,
        array $templateData,
    ): void {
        //do nothing
    }
}