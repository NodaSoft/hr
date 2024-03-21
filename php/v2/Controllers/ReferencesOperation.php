<?php

namespace NW\WebService\References\Operations\Notification\Controllers;

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName)
    {
        return $_REQUEST[$pName];
    }
}
