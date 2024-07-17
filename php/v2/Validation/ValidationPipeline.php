<?php

namespace NW\WebService\References\Operations\Notification\Validation;

use NW\WebService\References\Operations\Notification\Dto\NotificationData;
use NW\WebService\References\Operations\Notification\Notification\Exceptions\ValidationException;

class ValidationPipeline
{
    /**
     * @var ValidatorInterface[]
     */
    private array $validators = [];

    public function addValidator(ValidatorInterface $validator): self
    {
        $this->validators[] = $validator;
        return $this;
    }

    /**
     * @throws ValidationException
     */
    public function process(NotificationData $data): void
    {
        foreach ($this->validators as $validator) {
            $validator->validate($data);
        }
    }
}