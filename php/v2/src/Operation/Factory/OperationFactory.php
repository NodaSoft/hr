<?php

namespace NodaSoft\Operation\Factory;

use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Operation\FetchInitialData\FetchInitialData;
use NodaSoft\Operation\Command\Command;

interface OperationFactory
{
    public function setDependencies(Dependencies $dependencies): void;

    public function getFetchInitialData(): FetchInitialData;

    public function getCommand(): Command;
}
