<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Domain;

use Exception;
use NW\WebService\References\Operations\Notification\Struct\Request;
use NW\WebService\References\Operations\Notification\Struct\Result;

abstract class ReferencesOperation
{
    abstract public function doOperation(): Result;

    /**
     * @param non-empty-string $pName
     * @throws Exception
     */
    final public function getRequest(string $pName): Request
    {
        return new Request(...$_REQUEST[$pName]);
    }
}
