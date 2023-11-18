<?php

namespace NodaSoft\Operation\Factory;

use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Operation\Command\NotifyComplaintStatusChangedCommand;
use NodaSoft\Operation\FetchInitialData\FetchInitialData;
use NodaSoft\Operation\FetchInitialData\NotifyComplaintStatusChangedFetchInitialData;
use NodaSoft\Operation\Command\Command;

class NotifyComplaintStatusChangedFactory implements OperationFactory
{
    /** @var Dependencies */
    private $dependencies;

    public function setDependencies(Dependencies $dependencies): void
    {
        $this->dependencies = $dependencies;
    }
    /**
     * @return NotifyComplaintStatusChangedFetchInitialData
     */
    public function getFetchInitialData(): FetchInitialData {
        $fetch = new NotifyComplaintStatusChangedFetchInitialData();
        $fetch->setMapperFactory($this->dependencies->getMapperFactory());
        return $fetch;
    }

    /**
     * @return NotifyComplaintStatusChangedCommand
     */
    public function getCommand(): Command
    {
        $command = new NotifyComplaintStatusChangedCommand();
        $command->setEmail($this->dependencies->getEmailService());
        $command->setSms($this->dependencies->getSmsService());
        return $command;
    }
}
