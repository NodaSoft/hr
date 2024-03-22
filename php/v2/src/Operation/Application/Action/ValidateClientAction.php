<?php

namespace Src\Operation\Application\Action;

use Src\Operation\Application\DataTransferObject\ContractorData;
use Src\Operation\Application\DataTransferObject\SellerData;
use src\Operation\Application\Exceptions\ClientNotFoundException;
use src\Operation\Infrastructure\Domain\Enum\ContractorType;

class ValidateClientAction
{

    /**
     * @throws ClientNotFoundException
     */
    public function execute(ContractorData $client, SellerData $sellerData): void
    {
        if ($client->type !== ContractorType::TYPE_CUSTOMER || $client->sellerId !== $sellerData->id) {
            throw new ClientNotFoundException('Client not found', 400);
        }
    }
}