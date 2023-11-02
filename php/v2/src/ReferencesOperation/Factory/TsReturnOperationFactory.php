<?php

namespace NodaSoft\ReferencesOperation\Factory;

use NodaSoft\OperationParams\ReferencesOperationParams;
use NodaSoft\OperationParams\TsReturnOperationParams;
use NodaSoft\ReferencesOperation\Command\ReferencesOperationCommand;
use NodaSoft\ReferencesOperation\Command\TsReturnOperationCommand;
use NodaSoft\Request\Request;
use NodaSoft\Result\Operation\ReferencesOperation\ReferencesOperationResult;
use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;

class TsReturnOperationFactory implements ReferencesOperationFactory
{
    /** @var Request */
    private $request;

    public function setRequest(Request $request): void
    {
        $this->request = $request;
    }

    /**
     * @return TsReturnOperationResult
     */
    public function getResult(): ReferencesOperationResult
    {
        return new TsReturnOperationResult();
    }

    /**
     * @return TsReturnOperationParams
     */
    public function getParams(): ReferencesOperationParams
    {
        $params = new TsReturnOperationParams();
        $params->setRequest($this->request);
        return $params;
    }

    public function getCommand(
        ReferencesOperationResult $result,
        ReferencesOperationParams $params
    ): ReferencesOperationCommand
    {
        $command = new TsReturnOperationCommand();
        $command->setResult($result);
        $command->setParams($params);
        return $command;
    }
}
