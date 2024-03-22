<?php

namespace Src\Operation\Infrastructure\Adapters;

use Src\Contractor\Infrastructure\API\ContractorApi;
use Src\Operation\Application\DataTransferObject\SellerData;
use src\Operation\Application\Exceptions\SellerNotFoundException;

readonly class SellerAdapter
{
    private ContractorApi $employeeApi;

    public function __construct()
    {
        $this->employeeApi = new ContractorApi();
    }

    /**
     * @throws SellerNotFoundException
     */
    public function getById(int $contractorId): SellerData
    {
        $seller = $this->employeeApi->getById($contractorId);

        if ($seller == null) {
            throw new SellerNotFoundException('Seller not found', 400);
        }

        return SellerData::fromArray($seller);
    }

}