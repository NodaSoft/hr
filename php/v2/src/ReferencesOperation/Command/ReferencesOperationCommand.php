<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;

interface ReferencesOperationCommand
{
    public function execute(): ReferencesOperationResult;

    public function setResult(ReferencesOperationResult $result): void;

    public function setInitialData(InitialData $initialData): void;
}
