<?php

namespace NW\WebService\References\Operations\Notification\Services\Operations;

/**
 * Class ReferencesOperation
 * @package NW\WebService\References\Operations\Notification\Services
 */
abstract class ReferencesOperation
{
    /**
     * @return array
     */
    abstract public function doOperation(): array;

    /**
     * @param string $pName
     * @return array
     */
    public function getRequest(string $pName): array
    {
        return $_REQUEST[$pName] ?? [];
    }
}
