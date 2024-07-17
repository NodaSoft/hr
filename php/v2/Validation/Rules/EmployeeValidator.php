<?php

namespace NW\WebService\References\Operations\Notification\Validation;

use NW\WebService\References\Operations\Notification\Dto\NotificationData;
use NW\WebService\References\Operations\Notification\Employee;
use NW\WebService\References\Operations\Notification\Notification\Exceptions\ValidationException;

/**
 * Валидатор для проверки сотрудников
 */
class EmployeeValidator implements ValidatorInterface
{
    public function __construct(private readonly Employee $employee)
    {
    }

    public function validate(NotificationData $data): void
    {
        if ($this->employee->getById($data->creatorId) === null) {
            throw new ValidationException('Creator not found!');
        }
        if ($this->employee->getById($data->expertId) === null) {
            throw new ValidationException('Expert not found!');
        }
    }
}