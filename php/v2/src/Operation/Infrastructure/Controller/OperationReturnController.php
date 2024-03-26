<?php
declare(strict_types=1);

namespace src\Operation\Infrastructure\Controller;

use Src\Operation\Application\DataTransferObject\OperationData;
use src\Operation\Application\Exceptions\ClientNotFoundException;
use src\Operation\Application\Exceptions\ContractorNotFoundException;
use src\Operation\Application\Exceptions\EmployeeNotFoundException;
use src\Operation\Application\Exceptions\SellerNotFoundException;
use src\Operation\Application\Exceptions\ValidationException;
use src\Operation\Application\Request\OperationRequest;
use Src\Operation\Application\Service\OperationService;

final class OperationReturnController
{
    /**
     * @return array
     * @throws ValidationException
     * @throws ClientNotFoundException
     * @throws ContractorNotFoundException
     * @throws EmployeeNotFoundException
     * @throws SellerNotFoundException
     */
    public function httpRequestHandler(): array
    {
        $request = (new OperationRequest($_POST('data')))->validated();
        $operationData = OperationData::fromArray($request);
        $service = new OperationService();

        return $service->sendReturnNotification($operationData);
    }

}