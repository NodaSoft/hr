<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

abstract class AbstractReferencesOperation
{
    abstract public function doOperation(): array;

    protected function getRequest($pName): mixed
    {
        return $_REQUEST[$pName] ?? null;
    }
}
