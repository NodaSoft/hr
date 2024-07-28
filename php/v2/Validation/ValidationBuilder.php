<?php

namespace NW\WebService\References\Operations\Notification\Validation;

use Exception;

class ValidationBuilder
{
    /**
     * @var ValidatorInterface[]
     */
    private array $validators = [];

    /**
     * Adds a validator to the chain.
     *
     * @param ValidatorInterface $validator
     * @return self
     */
    public function add(ValidatorInterface $validator): self
    {
        $this->validators[] = $validator;
        return $this;
    }

    /**
     * Executes the validation chain.
     *
     * @param array $data
     * @param array $result
     * @return bool
     * @throws Exception
     */
    public function validate(array $data, array &$result): bool
    {
        foreach ($this->validators as $validator) {
            if (!$validator->validate($data, $result)) {
                return false;
            }
        }
        return true;
    }
}
