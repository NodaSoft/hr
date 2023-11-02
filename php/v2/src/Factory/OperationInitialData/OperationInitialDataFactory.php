<?php

namespace NodaSoft\Factory\OperationInitialData;

use NodaSoft\OperationInitialData\OperationInitialData;
use NodaSoft\OperationParams\ReferencesOperationParams;

interface OperationInitialDataFactory
{
    public function makeInitialData(ReferencesOperationParams $params): OperationInitialData;
}
