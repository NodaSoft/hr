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
use NodaSoft\Request\Request;

class NotifyComplaintStatusChangedFactory implements OperationFactory
{
    /** @var Request */
    private $request;

    public function setRequest(Request $request): void
    {
        $this->request = $request;
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
     * @param NotifyComplaintStatusChangedInitialData $initialData
     * @return NotifyComplaintStatusChangedCommand
     */
    public function getCommand(
        InitialData  $initialData,
        Dependencies $dependencies
    ): Command
    {
        $command = new NotifyComplaintStatusChangedCommand();
        $command->setInitialData($initialData);
        $command->setMail($dependencies->getMailService());
        $command->setSms($dependencies->getSmsService());
        return $command;
    }
}
