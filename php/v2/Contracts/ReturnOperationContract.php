<?php

namespace NW\WebService\References\Operations\Notification\Contracts;

use NW\WebService\References\Operations\Notification\Exceptions\ReturnOperationException;

interface ReturnOperationContract
{
    /**
     * @return array{
     *      notificationEmployeeByEmail: bool,
     *      notificationClientByEmail: bool,
     *      notificationClientBySms: array{isSent: bool, message: string},
     *  }
     *
     * @throws ReturnOperationException
     */
    public function doOperation(): array;
}
