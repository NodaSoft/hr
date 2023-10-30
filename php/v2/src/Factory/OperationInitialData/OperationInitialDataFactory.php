<?php

namespace NodaSoft\Factory\OperationInitialData;

use NodaSoft\OperationInitialData\OperationInitialData;
use NodaSoft\OperationParams\OperationParams;

interface OperationInitialDataFactory
{
    public function makeInitialData(OperationParams $params): OperationInitialData;
}
