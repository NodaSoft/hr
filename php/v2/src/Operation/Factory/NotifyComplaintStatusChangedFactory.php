<?php

namespace NodaSoft\Operation\Factory;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Operation\Command\NotifyComplaintStatusChangedCommand;
use NodaSoft\Operation\FetchInitialData\FetchInitialData;
use NodaSoft\Operation\FetchInitialData\NotifyComplaintStatusChangedFetchInitialData;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\InitialData\NotifyComplaintStatusChangedInitialData;
use NodaSoft\Operation\Params\Params;
use NodaSoft\Operation\Command\Command;
use NodaSoft\Operation\Params\NotifyComplaintStatusChangedParams;
use NodaSoft\Operation\Result\ReturnOperationStatusChangedResult;
use NodaSoft\Request\Request;
use NodaSoft\Operation\Result\Result;

class NotifyComplaintStatusChangedFactory implements OperationFactory
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
    public function getResult(): Result
    {
        return new ReturnOperationStatusChangedResult();
    }

    /**
     * @return NotifyComplaintStatusChangedParams
     */
    public function getParams(): Params
    {
        $params = new NotifyComplaintStatusChangedParams();
        $params->setRequest($this->request);
        return $params;
    }

    /**
     * @param MapperFactory $mapperFactory
     * @return NotifyComplaintStatusChangedFetchInitialData
     */
    public function getFetchInitialData(
        MapperFactory $mapperFactory
    ): FetchInitialData {
        $fetch = new NotifyComplaintStatusChangedFetchInitialData();
        $fetch->setMapperFactory($mapperFactory);
        return $fetch;
    }

    /**
     * @param Result $result
     * @param NotifyComplaintStatusChangedInitialData $initialData
     * @return NotifyComplaintStatusChangedCommand
     */
    public function getCommand(
        Result       $result,
        InitialData  $initialData,
        Dependencies $dependencies
    ): Command
    {
        $command = new NotifyComplaintStatusChangedCommand();
        $command->setResult($result);
        $command->setInitialData($initialData);
        $command->setMail($dependencies->getMailService());
        $command->setSms($dependencies->getSmsService());
        return $command;
    }
}
