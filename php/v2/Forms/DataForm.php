<?php

namespace NW\WebService\References\Operations\Notification\Forms;

use NW\WebService\References\Operations\Notification\Helpers\Form;
use NW\WebService\References\Operations\Notification\Exceptions\ValidationException;
use NW\WebService\References\Operations\Notification\models\Contractor;
use NW\WebService\References\Operations\Notification\models\Employee;
use NW\WebService\References\Operations\Notification\models\Seller;

// Обработка и валидация входных данных
class DataForm extends Form
{
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
     * @var int
     */
    public mixed $complaintId = null;
    /**
     * @var string
     */
    public mixed $complaintNumber = null;
    /**
     * @var int
     */
    public mixed $consumptionId = null;
    /**
     * @var string
     */
    public mixed $consumptionNumber = null;
    /**
     * @var string
     */
    public mixed $agreementNumber = null;
    /**
     * @var string
     */
    public mixed $date = null;
    /**
     * @var Seller
     */
    public mixed $reseller;
    /**
     * @var Contractor
     */
    public mixed $client;
    /**
     * @var Employee
     */
    public mixed $creator;
    /**
     * @var Employee
     */
    public mixed $expert;

    /**
     * @var array
     */
    public mixed $differences;

    private const NOT_FOUND_EXCEPTION = '%s not found!';
    public function rules(): array
    {
        return [
            'resellerId' => ['required', 'integer', function (string $attribute) {
                $this->reseller = $this->getReseller($this->$attribute);
                if ($this->reseller === null) {
                    $this->addError($attribute, sprintf(self::NOT_FOUND_EXCEPTION, $attribute));
                    return false;
                }
                return true;
            }],
            'clientId' => ['required', 'integer', function (string $attribute) {
                $this->client = $this->getClient($this->$attribute);
                if ($this->client === null) {
                    $this->addError($attribute, sprintf(self::NOT_FOUND_EXCEPTION, $attribute));
                    return false;
                }
                return true;
            }],
            'creatorId' => ['required', 'integer', function (string $attribute) {
                $this->creator = $this->getEmployee($this->$attribute);
                if ($this->creator === null) {
                    $this->addError($attribute, sprintf(self::NOT_FOUND_EXCEPTION, $attribute));
                    return false;
                }
                return true;
            }],
            'expertId' => ['required', 'integer', function (string $attribute) {
                $this->expert = $this->getEmployee($this->$attribute);
                if ($this->expert === null) {
                    $this->addError($attribute, sprintf(self::NOT_FOUND_EXCEPTION, $attribute));
                    return false;
                }
                return true;
            }],
            'notificationType' => ['required', 'integer'],
            'differences' => [
                function (string $attribute) {
                    if (!is_array($this->$attribute)) {
                        $this->$attribute = []; // или добавить валидацию
                    }
                }
            ],
            'complaintNumber' => ['required', 'string'],
            'consumptionId' => ['required', 'integer'],
            'consumptionNumber' => ['required', 'string'],
            'agreementNumber' => ['required', 'string'],
            'date' => ['required', 'string'],
        ];
    }

    private function getReseller(int $resellerId): ?Seller
    {
        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            return null;
        }
        return $reseller;
    }

    private function getClient(int $clientId): ?Contractor
    {
        $client = Contractor::getById($clientId);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $this->resellerId) {
            return null;
        }
        return $client;
    }

    private function getEmployee(int $employeeId): ?Employee
    {
        $employee = Employee::getById($employeeId);
        if ($employee === null) {
            return null;
        }
        return $employee;
    }
}