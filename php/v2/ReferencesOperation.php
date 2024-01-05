<?php

namespace NW\WebService\References\Operations\Notification;

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    /**
     * @param string $pName
     * @return mixed
     */
    public function getRequest(string $pName)
    {
        if (!key_exists($pName, $_REQUEST)) {
            return null;
        }

        return $_REQUEST[$pName];
    }
}
