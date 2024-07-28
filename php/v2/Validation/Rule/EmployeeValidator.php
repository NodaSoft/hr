<?php

namespace NW\WebService\References\Operations\Notification\Validation\Rule;

use NW\WebService\References\Operations\Notification\Employee;
use NW\WebService\References\Operations\Notification\Validation\ValidatorInterface;

/**
 * EmployeeValidator class
 */
class EmployeeValidator implements ValidatorInterface
{

    public function validate(array $data, array &$result = []): bool
    {
        $creatorId = (int)($data['creatorId'] ?? 0);
        $expertId = (int)($data['expertId'] ?? 0);

        $cr = Employee::getById($creatorId);
        if (empty($cr)) {
            throw new \Exception('Creator not found!', 400);
        }

        $et = Employee::getById($expertId);
        if (empty($et)) {
            throw new \Exception('Expert not found!', 400);
        }

        return true;
    }
}
