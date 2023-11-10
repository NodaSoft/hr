<?php

namespace NodaSoft\ReferencesOperation\Factory;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\ReferencesOperation\FetchInitialData\FetchInitialData;
use NodaSoft\ReferencesOperation\FetchInitialData\TsReturnFetchInitialData;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\TsReturnInitialData;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Params\TsReturnOperationParams;
use NodaSoft\ReferencesOperation\Command\ReferencesOperationCommand;
use NodaSoft\ReferencesOperation\Command\TsReturnOperationCommand;
use NodaSoft\Request\Request;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;
use NodaSoft\ReferencesOperation\Result\TsReturnOperationResult;

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

    /**
     * @param MapperFactory $mapperFactory
     * @return FetchInitialData
     */
    public function getFetchInitialData(
        MapperFactory $mapperFactory
    ): FetchInitialData {
        $fetch = new TsReturnFetchInitialData();
        $fetch->setMapperFactory($mapperFactory);
        return $fetch;
    }

    /**
     * @param ReferencesOperationResult $result
     * @param TsReturnInitialData $initialData
     * @return ReferencesOperationCommand
     */
    public function getCommand(
        ReferencesOperationResult $result,
        InitialData $initialData,
        Dependencies $dependencies
    ): ReferencesOperationCommand
    {
        $command = new TsReturnOperationCommand();
        $command->setResult($result);
        $command->setInitialData($initialData);
        $command->setMail($dependencies->getMailService());
        $command->setSms($dependencies->getSmsService());
        return $command;
    }
}
