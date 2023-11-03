<?php

namespace NodaSoft\Factory\OperationInitialData;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\OperationInitialData\OperationInitialData;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;

interface OperationInitialDataFactory
{
    public function setMapperFactory(MapperFactory $mapperFactory): void;

    public function makeInitialData(ReferencesOperationParams $params): OperationInitialData;
}
