<?php

namespace NW\WebService\References\Operations\Notification\Forms;

use NW\WebService\References\Operations\Notification\Exceptions\ValidationException;
use NW\WebService\References\Operations\Notification\models\Contractor;
use NW\WebService\References\Operations\Notification\models\Employee;
use NW\WebService\References\Operations\Notification\models\Seller;

// Обработка и валидация входных данных
class IndexForm
{
    private const RESELLER_ID_FIELD = 'resellerId';
    private const CLIENT_ID_FIELD = 'clientId';
    private const EXPERT_ID_FIELD = 'expertId';
    private const CREATOR_ID_FIELD = 'creatorId';
    private const NOTIFICATION_TYPE_FIELD = 'notificationType';
    private const DIFFERENCES_FIELD = 'differences';

    private const ERROR_CLIENT_NOT_FOUND = 'Client not found!';
    private const ERROR_CREATOR_NOT_FOUND = 'Creator not found!';
    private const ERROR_EXPERT_NOT_FOUND = 'Expert not found!';
    private const ERROR_SELLER_NOT_FOUND = 'Seller not found!';
    private const FORMAT_EMPTY_MESSAGE = 'Empty %s';
    /**
     * @var int
     */
    public mixed $resellerId = null;
    /**
     * @var int
     */
    public mixed $clientId = null;
    /**
     * @var int
     */
    public mixed $creatorId = null;
    /**
     * @var int
     */
    public mixed $expertId = null;
    /**
     * @var int
     */
    public mixed $notificationType = null;
    /**
     * @var Seller
     */
    public Seller $reseller;
    /**
     * @var Contractor
     */
    public Contractor $client;
    /**
     * @var Employee
     */
    public Employee $creator;
    /**
     * @var Employee
     */
    public Employee $expert;

    /**
     * @var array
     */
    public array $differences;

    /**
     * @var array validation errors (attribute name => array of errors)
     */
    private array $_errors;

    /**
     * @throws ValidationException
     */
    public function __construct(
        mixed $data
    ) {
        $this->resellerId = $this->checkRequiredField(self::RESELLER_ID_FIELD, $data);
        $this->clientId = $this->checkRequiredField(self::CLIENT_ID_FIELD, $data);
        $this->creatorId = $this->checkRequiredField(self::CREATOR_ID_FIELD, $data);
        $this->expertId = $this->checkRequiredField(self::EXPERT_ID_FIELD, $data);
        $this->notificationType = $this->checkRequiredField(self::NOTIFICATION_TYPE_FIELD, $data);
        $this->differences = $data[self::DIFFERENCES_FIELD] ?? [];

        $reseller = $this->getReseller();

        $client = $this->getClient();

        $creator = $this->getEmployee($this->creatorId, self::CREATOR_ID_FIELD, self::ERROR_CREATOR_NOT_FOUND);

        $expert = $this->getEmployee($this->expertId, self::EXPERT_ID_FIELD, self::ERROR_EXPERT_NOT_FOUND);

        if (count($this->_errors) > 0) {
            throw new ValidationException($this->_errors);
        }
        $this->reseller = $reseller;
        $this->client = $client;
        $this->creator = $creator;
        $this->expert = $expert;
    }

    private function getReseller(): ?Seller
    {
        $reseller = Seller::getById($this->resellerId);
        if ($reseller === null) {
            $this->addError(self::RESELLER_ID_FIELD, self::ERROR_SELLER_NOT_FOUND);
            return null;
        }
        return $reseller;
    }

    private function getClient(): ?Contractor
    {
        $client = Contractor::getById($this->clientId);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $this->resellerId) {
            $this->addError(self::CLIENT_ID_FIELD, self::ERROR_CLIENT_NOT_FOUND);
            return null;
        }
        return $client;
    }

    private function getEmployee(int $employeeId, string $errorField, string $errorMessage): ?Employee
    {
        $employee = Employee::getById($employeeId);
        if ($employee === null) {
            $this->addError($errorField, $errorMessage);
            return null;
        }
        return $employee;
    }

    public function addError($attribute, $message): void
    {
        $this->_errors[$attribute][] = $message;
    }

    public function generateErrorMessage(string $format, string $attribute): string
    {
        return sprintf($format, $attribute);
    }

    public function checkRequiredField(string $field, mixed $data): int
    {
        $value = (int)$data[$field];
        if ($value === 0) {
            $this->addError($field, $this->generateErrorMessage(self::FORMAT_EMPTY_MESSAGE, $field));
        }
        return $value;
    }
}