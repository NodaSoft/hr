<?php

namespace NodaSoft\Operation\FetchInitialData;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Request\Request;

interface FetchInitialData
{
    public function setMapperFactory(MapperFactory $mapperFactory): void;

    public function fetch(Request $request): InitialData;
}
