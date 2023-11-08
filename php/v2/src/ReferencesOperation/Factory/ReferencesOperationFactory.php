<?php

namespace NodaSoft\ReferencesOperation\Factory;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\ReferencesOperation\FetchInitialData\FetchInitialData;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Command\ReferencesOperationCommand;
use NodaSoft\Request\Request;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;

interface ReferencesOperationFactory
{
    public function setRequest(Request $request): void;

    public function getResult(): ReferencesOperationResult;

    public function getParams(): ReferencesOperationParams;

    public function getFetchInitialData(
        MapperFactory $mapperFactory
    ): FetchInitialData;

    /**
     * @param ReferencesOperationResult $result
     * @param InitialData $initialData
     * @param Dependencies $dependencies
     * @return ReferencesOperationCommand
     */
    public function getCommand(
        ReferencesOperationResult $result,
        InitialData $initialData,
        Dependencies $dependencies
    ): ReferencesOperationCommand;
}
