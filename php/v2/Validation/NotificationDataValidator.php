<?php

namespace NW\WebService\References\Operations\Notification\Validation;

use NW\WebService\References\Operations\Notification\Contractor;
use NW\WebService\References\Operations\Notification\Employee;
use NW\WebService\References\Operations\Notification\Seller;

/**
 * Реализация валидатора данных уведомления
 */
readonly class NotificationDataValidator
{

    public function __construct(private Seller $seller, private Contractor $contractor, private Employee $employee)
    {
    }

    public function isValid(): ValidationPipeline
    {
        // Создаем валидаторы
        $sellerValidator = new SellerValidator($this->seller);
        $clientValidator = new ClientValidator($this->contractor);
        $employeeValidator = new EmployeeValidator($this->employee);

        // Создаем и настраиваем пайплайн валидации
        $validationPipeline = new ValidationPipeline();
        $validationPipeline
            ->addValidator($sellerValidator)
            ->addValidator($clientValidator)
            ->addValidator($employeeValidator);

        return $validationPipeline;
    }
}