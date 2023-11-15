<?php

namespace NodaSoft\Operation\Command;

use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\Result\Result;

interface Command
{
    public function execute(): Result;

    public function setResult(Result $result): void;

    public function setInitialData(InitialData $initialData): void;
}
