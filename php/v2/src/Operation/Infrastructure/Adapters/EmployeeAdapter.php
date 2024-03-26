<?php

namespace Src\Operation\Infrastructure\Adapters;

use Src\Employee\Infrastructure\API\EmployeeApi;
use Src\Operation\Application\DataTransferObject\EmployeeData;
use src\Operation\Application\Exceptions\EmployeeNotFoundException;

readonly class EmployeeAdapter
{
    private EmployeeApi $employeeApi;

    public function __construct()
    {
        $this->employeeApi = new EmployeeApi();
    }

    /**
     * @throws EmployeeNotFoundException
     */
    public function getById(int $employeeId): EmployeeData
    {
        $employee = $this->employeeApi->getById($employeeId);

        if ($employee == null) {
            throw new EmployeeNotFoundException('Seller not found', 400);
        }
        return EmployeeData::fromArray($employee);
    }

}