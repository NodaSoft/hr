<?php

namespace app\Http\Controllers\References\Operations;

abstract class OperationController
{
    abstract public function doOperation(): array;

    public function getRequest($pName)
    {
        return $_REQUEST[$pName];
    }
}