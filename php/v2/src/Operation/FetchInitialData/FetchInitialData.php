<?php

namespace NodaSoft\Operation\FetchInitialData;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\Params\Params;

interface FetchInitialData
{
    public function setMapperFactory(MapperFactory $mapperFactory): void;

    public function fetch(Params $params): InitialData;
}
