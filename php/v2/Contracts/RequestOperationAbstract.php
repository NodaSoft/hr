<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Contracts;

/**
 * Class RequestOperationAbstarct.
 *
 */
abstract class RequestOperationAbstract
{
    /**
     * get Status name by id
     *
     * @return array
     */
    abstract public function doOperation(): array;

    /**
     * get Status name by id
     *
     * @param $pName
     * @return mixed|null
     */
    public function getRequest($pName)
    {
        return $_REQUEST[$pName] ?? null;
    }
}