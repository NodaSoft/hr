<?php

namespace NodaSoft\Operation\Factory;

use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Operation\FetchInitialData\FetchInitialData;
use NodaSoft\Operation\FetchInitialData\NotifyComplaintNewFetchInitialData;
use NodaSoft\Operation\Command\Command;
use NodaSoft\Operation\Command\NotifyComplaintNewCommand;

class NotifyComplaintNewFactory implements OperationFactory
{
    /** @var Dependencies */
    private $dependencies;

    public function setDependencies(Dependencies $dependencies): void
    {
        $this->dependencies = $dependencies;
    }

    /**
     * @return NotifyComplaintNewFetchInitialData
     */
    public function getFetchInitialData(): FetchInitialData {
        $fetch = new NotifyComplaintNewFetchInitialData();
        $fetch->setMapperFactory($this->dependencies->getMapperFactory());
        return $fetch;
    }

    /**
     * @return NotifyComplaintNewCommand
     */
    public function getCommand(): Command
    {
        $command = new NotifyComplaintNewCommand();
        $command->setEmail($this->dependencies->getEmailService());
        return $command;
    }
}
