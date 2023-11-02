<?php

namespace NodaSoft\ReferencesOperation\Factory;

use NodaSoft\OperationParams\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Command\ReferencesOperationCommand;
use NodaSoft\Request\Request;
use NodaSoft\Result\Operation\ReferencesOperation\ReferencesOperationResult;

interface ReferencesOperationFactory
{
    public function setRequest(Request $request): void;

    public function getResult(): ReferencesOperationResult;

    public function getParams(): ReferencesOperationParams;

    public function getCommand(
        ReferencesOperationResult $result,
        ReferencesOperationParams $params
    ): ReferencesOperationCommand;
}
