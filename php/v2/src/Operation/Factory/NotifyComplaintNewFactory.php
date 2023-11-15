<?php

namespace NodaSoft\Operation\Factory;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Operation\FetchInitialData\FetchInitialData;
use NodaSoft\Operation\FetchInitialData\NotifyComplaintNewFetchInitialData;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\InitialData\NotifyComplaintNewInitialData;
use NodaSoft\Operation\Params\Params;
use NodaSoft\Operation\Params\NotifyComplaintNewParams;
use NodaSoft\Operation\Command\Command;
use NodaSoft\Operation\Command\NotifyComplaintNewCommand;
use NodaSoft\Request\Request;
use NodaSoft\Operation\Result\Result;
use NodaSoft\Operation\Result\NotifyComplaintNewResult;

class NotifyComplaintNewFactory implements OperationFactory
{
    /** @var Request */
    private $request;

    public function setRequest(Request $request): void
    {
        $this->request = $request;
    }

    /**
     * @return NotifyComplaintNewResult
     */
    public function getResult(): Result
    {
        return new NotifyComplaintNewResult();
    }

    /**
     * @return NotifyComplaintNewParams
     */
    public function getParams(): Params
    {
        $params = new NotifyComplaintNewParams();
        $params->setRequest($this->request);
        return $params;
    }

    /**
     * @param MapperFactory $mapperFactory
     * @return NotifyComplaintNewFetchInitialData
     */
    public function getFetchInitialData(
        MapperFactory $mapperFactory
    ): FetchInitialData {
        $fetch = new NotifyComplaintNewFetchInitialData();
        $fetch->setMapperFactory($mapperFactory);
        return $fetch;
    }

    /**
     * @param NotifyComplaintNewResult $result
     * @param NotifyComplaintNewInitialData $initialData
     * @return NotifyComplaintNewCommand
     */
    public function getCommand(
        Result       $result,
        InitialData  $initialData,
        Dependencies $dependencies
    ): Command
    {
        $command = new NotifyComplaintNewCommand();
        $command->setResult($result);
        $command->setInitialData($initialData);
        $command->setMail($dependencies->getMailService());
        return $command;
    }
}
