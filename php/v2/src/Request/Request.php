<?php

namespace NodaSoft\Request;

use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Operation\Factory\OperationFactory;

interface Request
{
    /**
     * @param string $key
     * @return mixed
     */
    public function get(string $key);

    public function getOperationFactory(
        Dependencies $dependencies
    ): OperationFactory;
}
