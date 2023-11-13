<?php

namespace NodaSoft\ReferencesOperation\Factory;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\ReferencesOperation\Command\ReturnOperationStatusChangedCommand;
use NodaSoft\ReferencesOperation\FetchInitialData\FetchInitialData;
use NodaSoft\ReferencesOperation\FetchInitialData\ReturnOperationStatusChangedFetchInitialData;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationStatusChangedInitialData;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Command\ReferencesOperationCommand;
use NodaSoft\ReferencesOperation\Params\ReturnOperationStatusChangedParams;
use NodaSoft\ReferencesOperation\Result\ReturnOperationStatusChangedResult;
use NodaSoft\Request\Request;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;

class ReturnOperationStatusChangedFactory implements ReferencesOperationFactory
{
    /** @var Request */
    private $request;

    public function setRequest(Request $request): void
    {
        $this->request = $request;
    }

    /**
     * @return ReturnOperationStatusChangedResult
     */
    public function getResult(): ReferencesOperationResult
    {
        return new ReturnOperationStatusChangedResult();
    }

    /**
     * @return ReturnOperationStatusChangedParams
     */
    public function getParams(): ReferencesOperationParams
    {
        $params = new ReturnOperationStatusChangedParams();
        $params->setRequest($this->request);
        return $params;
    }

    /**
     * @param MapperFactory $mapperFactory
     * @return ReturnOperationStatusChangedFetchInitialData
     */
    public function getFetchInitialData(
        MapperFactory $mapperFactory
    ): FetchInitialData {
        $fetch = new ReturnOperationStatusChangedFetchInitialData();
        $fetch->setMapperFactory($mapperFactory);
        return $fetch;
    }

    /**
     * @param ReferencesOperationResult $result
     * @param ReturnOperationStatusChangedInitialData $initialData
     * @return ReturnOperationStatusChangedCommand
     */
    public function getCommand(
        ReferencesOperationResult $result,
        InitialData $initialData,
        Dependencies $dependencies
    ): ReferencesOperationCommand
    {
        $command = new ReturnOperationStatusChangedCommand();
        $command->setResult($result);
        $command->setInitialData($initialData);
        $command->setMail($dependencies->getMailService());
        $command->setSms($dependencies->getSmsService());
        return $command;
    }
}
