<?php

namespace NW\WebService\References\Operations\Notification\Contracts;

use NW\WebService\References\Operations\Notification\Contractor;

interface ContractorServiceContract
{
    public function getById(int $id): ?Contractor;
}
