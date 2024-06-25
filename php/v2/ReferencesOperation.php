<?php

namespace NW\WebService\References\Operations\Notification;

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    public function getRequest($pName): array
    {
        return (array) $_REQUEST[$pName];
    }
}
