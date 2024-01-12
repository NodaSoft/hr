<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    /**
     * @param $pName
     * @return mixed|null
     */
    public function getRequest($pName): mixed
    {
        return $_REQUEST[$pName] ?? null;
    }
}
