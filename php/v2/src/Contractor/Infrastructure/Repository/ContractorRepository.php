<?php

namespace Src\Contractor\Infrastructure\Repository;

use Src\Contractor\Domain\Entity\Contractor;
use Src\Contractor\Domain\Repository\ContractorRepositoryInterface;

class ContractorRepository implements ContractorRepositoryInterface
{

    public function getById(int $contractorId): Contractor
    {
        // TODO: Implement getById() method.
    }
}