<?php

namespace NodaSoft\Factory\OperationInitialData;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\OperationInitialData\OperationInitialData;
use NodaSoft\OperationParams\ReferencesOperationParams;

interface OperationInitialDataFactory
{
    public function setMapperFactory(MapperFactory $mapperFactory): void;

    public function makeInitialData(ReferencesOperationParams $params): OperationInitialData;
}
