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

class NotifyComplaintNewFactory implements OperationFactory
{
    /** @var Request */
    private $request;

    public function setRequest(Request $request): void
    {
        $this->request = $request;
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
     * @param NotifyComplaintNewInitialData $initialData
     * @return NotifyComplaintNewCommand
     */
    public function getCommand(
        InitialData  $initialData,
        Dependencies $dependencies
    ): Command
    {
        $command = new NotifyComplaintNewCommand();
        $command->setInitialData($initialData);
        $command->setMail($dependencies->getMailService());
        return $command;
    }
}
