<?php

declare(strict_types=1);


namespace NW\WebService\Validation;

use NW\WebService\Employees\Employee;
use NW\WebService\Employees\EmployeeTypeEnum;
use NW\WebService\Exceptions\ValidationException;
use NW\WebService\Request\DTO\RequestDTO;

class ValidationRequest
{
    /**
     * @throws  ValidationException
     */
    public static function validate(RequestDTO $data): void
    {
        if ( ! $data->resellerId) {
            throw new ValidationException('Empty resellerId', 400);
        }

        if ( ! $data->clientId) {
            throw new ValidationException('Empty clientId', 400);
        }

        if ( ! $data->creatorId) {
            throw new ValidationException('Empty creatorId', 400);
        }

        if ( ! $data->expertId) {
            throw new ValidationException('Empty expertId', 400);
        }

        if ( ! Employee::getById(type: EmployeeTypeEnum::RESELLER, id: $data->resellerId)) {
            throw new ValidationException('reseller not found!', 400);
        }

        if ( ! Employee::getById(type: EmployeeTypeEnum::CONTRACTOR, id: $data->clientId)) {
            throw new ValidationException('client not found!', 400);
        }

        if ( ! Employee::getById(type: EmployeeTypeEnum::CREATOR, id: $data->creatorId)) {
            throw new ValidationException('Creator not found!', 400);
        }

        if ( ! Employee::getById(type: EmployeeTypeEnum::EXPERT, id: $data->expertId)) {
            throw new ValidationException('Expert not found!', 400);
        }
    }
}
