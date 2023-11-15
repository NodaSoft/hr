<?php

namespace NodaSoft\Operation\Factory;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Operation\FetchInitialData\FetchInitialData;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\Params\Params;
use NodaSoft\Operation\Command\Command;
use NodaSoft\Request\Request;
use NodaSoft\Operation\Result\Result;

interface OperationFactory
{
    public function setRequest(Request $request): void;

    public function getResult(): Result;

    public function getParams(): Params;

    public function getFetchInitialData(
        MapperFactory $mapperFactory
    ): FetchInitialData;

    /**
     * @param Result $result
     * @param InitialData $initialData
     * @param Dependencies $dependencies
     * @return Command
     */
    public function getCommand(
        Result       $result,
        InitialData  $initialData,
        Dependencies $dependencies
    ): Command;
}
