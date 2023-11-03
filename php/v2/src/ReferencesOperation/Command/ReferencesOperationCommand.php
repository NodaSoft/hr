<?php

namespace NodaSoft\ReferencesOperation\Command;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;

interface ReferencesOperationCommand
{
    public function execute(): ReferencesOperationResult;

    public function setResult(ReferencesOperationResult $result): void;

    public function setParams(ReferencesOperationParams $params): void;

    public function setMapperFactory(MapperFactory $mapperFactory): void;
}
