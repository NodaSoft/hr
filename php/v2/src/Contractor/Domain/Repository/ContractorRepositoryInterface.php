<?php

namespace Src\Contractor\Domain\Repository;

use Src\Contractor\Domain\Entity\Contractor;

interface ContractorRepositoryInterface
{
    public function getById(int $contractorId): Contractor;
}