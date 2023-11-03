<?php

namespace NodaSoft\ReferencesOperation\FetchInitialData;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;

interface FetchInitialData
{
    public function setMapperFactory(MapperFactory $mapperFactory): void;

    public function fetch(ReferencesOperationParams $params): InitialData;
}
