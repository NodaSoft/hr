<?php

namespace NodaSoft\ReferencesOperation\Factory;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\ReferencesOperation\FetchInitialData\FetchInitialData;
use NodaSoft\ReferencesOperation\FetchInitialData\ReturnOperationNewFetchInitialData;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationNewInitialData;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Params\ReturnOperationNewParams;
use NodaSoft\ReferencesOperation\Command\ReferencesOperationCommand;
use NodaSoft\ReferencesOperation\Command\ReturnOperationNewCommand;
use NodaSoft\Request\Request;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;
use NodaSoft\ReferencesOperation\Result\ReturnOperationNewResult;

class ReturnOperationNewFactory implements ReferencesOperationFactory
{
    /** @var Request */
    private $request;

    public function setRequest(Request $request): void
    {
        $this->request = $request;
    }

    /**
     * @return ReturnOperationNewResult
     */
    public function getResult(): ReferencesOperationResult
    {
        return new ReturnOperationNewResult();
    }

    /**
     * @return ReturnOperationNewParams
     */
    public function getParams(): ReferencesOperationParams
    {
        $params = new ReturnOperationNewParams();
        $params->setRequest($this->request);
        return $params;
    }

    /**
     * @param MapperFactory $mapperFactory
     * @return ReturnOperationNewFetchInitialData
     */
    public function getFetchInitialData(
        MapperFactory $mapperFactory
    ): FetchInitialData {
        $fetch = new ReturnOperationNewFetchInitialData();
        $fetch->setMapperFactory($mapperFactory);
        return $fetch;
    }

    /**
     * @param ReturnOperationNewResult $result
     * @param ReturnOperationNewInitialData $initialData
     * @return ReturnOperationNewCommand
     */
    public function getCommand(
        ReferencesOperationResult $result,
        InitialData $initialData,
        Dependencies $dependencies
    ): ReferencesOperationCommand
    {
        $command = new ReturnOperationNewCommand();
        $command->setResult($result);
        $command->setInitialData($initialData);
        $command->setMail($dependencies->getMailService());
        return $command;
    }
}
