<?php

namespace NodaSoft\ReferencesOperation\Factory;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Command\ReferencesOperationCommand;
use NodaSoft\Request\Request;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;

interface ReferencesOperationFactory
{
    public function setRequest(Request $request): void;

    public function getResult(): ReferencesOperationResult;

    public function getParams(): ReferencesOperationParams;

    public function getCommand(
        ReferencesOperationResult $result,
        ReferencesOperationParams $params,
        MapperFactory $mapperFactory
    ): ReferencesOperationCommand;
}
