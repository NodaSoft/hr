<?php

namespace NW\WebService\References\Operations\Notification\Services;

use NW\WebService\References\Operations\Notification\Models\Contractor;

/**
 * Class OperationValidator
 * @package NW\WebService\References\Operations\Notification\Models\Contractor
 */
class OperationValidator
{
    /** @var array $data */
    private array $data;

    /** @var string|null $notificationMessage */
    private ?string $notificationMessage;

    /**
     * @param array $data
     */
    public function __construct(array $data)
    {
        $this->data = $data;
    }

    /**
     * @param array $data
     * @return self
     */
    public static function make(array $data): self
    {
        return (new self($data));
    }

    /**
     * @return string|null
     */
    public function getNotificationMessage(): ?string
    {
        return $this->notificationMessage;
    }

    /**
     * @param string $notificationMessage
     * @return void
     */
    private function setNotificationMessage(string $notificationMessage): void
    {
        $this->notificationMessage = $notificationMessage;
    }

    /**
     * @return string|null
     */
    public function validate(): ?string
    {
        $rules = [
            'resellerId' => 'required_with_notify|int|exists:Seller',
            'notificationType' => 'required|string',
            'clientId' => 'required|int|exists:Contractor,Client|isCorrectClient',
            'creatorId' => 'required|int|exists:Employee,Creator',
            'expertId' => 'required|int|exists:Employee,Expert',
            'differences' => 'array',
            'differences.from' => 'required|int',
            'differences.to' => 'required|int',
            'complaintId' => 'required|int',
            'complaintNumber' => 'required|string',
            'consumptionId' => 'required|int',
            'consumptionNumber' => 'required|string',
            'agreementNumber' => 'required|string',
            'date' => 'required|string',
        ];

        foreach ($rules as $field => $rule) {
            $datum = $this->getDatum($field);
            $fieldRules = explode('|', $rule);
            foreach ($fieldRules as $fieldRule) {
                if ($fieldRule === 'required' && empty($datum)) {
                    return 'Empty ' . $field;
                }
                if (str_contains('exists:', $fieldRule)) {
                    $validationParams = explode(',', str_replace('exists:', '', $fieldRule));
                    $model = $validationParams[0];
                    $entity = $validationParams[1];

                    if (!'\NW\WebService\References\Operations\Notification\Models\\' . $model::getById($datum)) {
                        return $entity . ' not found!';
                    }
                }
                if ($fieldRule === 'required_with_notify' && empty($datum)) {
                    $this->setNotificationMessage("Empty $field");
                }
                if ($fieldRule === 'isCorrectClient') {
                    $client = Contractor::getById((int)$datum);
                    if ($client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $this->getDatum('resellerId')) {
                        return 'Client not found!';
                    }
                }

                // Some more custom validations
                if ($fieldRule === 'array' && !is_array($datum)) {
                    return 'Value of ' . $field . ' is invalid, must be array.';
                }
                if ($fieldRule === 'string' && !is_string($datum)) {
                    return 'Value of ' . $field . ' is invalid, must be string.';
                }
                if ($fieldRule === 'int' && !is_int($datum)) {
                    return 'Value of ' . $field . ' is invalid, must be integer.';
                }
            }
        }

        return null;
    }

    /**
     * @param string $key
     * @param mixed|null $datum
     * @return array|mixed|null
     */
    private function getDatum(string $key, mixed $datum = null): mixed
    {
        if (is_null($datum)) {
            $datum = $this->data;
        }

        if (!str_contains($key, '.')) {
            return $datum[$key] ?? null;
        }

        $keys = explode('.', $key);

        foreach ($keys as $key) {
            $datum = $this->getDatum($key, $datum);
        }

        return $datum;
    }
}
