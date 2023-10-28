<?php

namespace NodaSoft\Factory\OperationInitialData;

use NodaSoft\OperationInitialData\OperationInitialData;

interface OperationInitialDataFactory
{
    public function makeInitialData(array $params): OperationInitialData;
}
