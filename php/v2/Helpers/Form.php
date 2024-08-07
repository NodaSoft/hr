<?php

namespace NW\WebService\References\Operations\Notification\Helpers;

abstract class Form
{
    protected array $errors = [];

    public function load(array $data): bool
    {
        return $this->setAttributes($data);
    }

    private function setAttributes(array $data): bool
    {
        $success = false;
        foreach ($data as $name => $value) {
            if (property_exists($this, $name)) {
                $this->$name = $value;
                $success = true;
            }
        }
        return $success;
    }

    public function validate(array $attributeNames = null): bool
    {
        $success = true;
        $rules = $this->rules();

        if ($attributeNames !== null) {
            foreach ($attributeNames as $name) {
                if (!$this->validateAttribute($name, $rules)) {
                    $success = false;
                }
            }
        } else {
            foreach (array_keys($rules) as $name) {
                if (!$this->validateAttribute($name, $rules)) {
                    $success = false;
                }
            }
        }

        return $success;
    }

    private function validateAttribute(string $name, array $rules): bool
    {
        $success = true;
        if (isset($rules[$name])) {
            foreach ($rules[$name] as $rule) {
                if (is_string($rule) || is_callable($rule)) {
                    $validator = $rule;
                    $options = [];
                } else {
                    [$validator, $options] = $rule;
                }
                if (is_callable($validator)) {
                    $validator = $rule;
                    $options = [];
                    $isValid = $validator($name, $options);
                } else {
                    $isValid = $this->$validator($name, $options);
                }
                if (!$isValid) {
                    $success = false;
                }
            }
        }
        return $success;
    }

    protected function rules(): array
    {
        return [];
    }

    protected function addError(string $attribute, string $message): void
    {
        $this->errors[$attribute][] = $message;
    }

    public function hasErrors(string $attribute = null): bool
    {
        if ($attribute === null) {
            return !empty($this->errors);
        }
        return isset($this->errors[$attribute]);
    }

    public function getErrors(string $attribute = null): array
    {
        if ($attribute === null) {
            return $this->errors;
        }
        return $this->errors[$attribute] ?? [];
    }

    public function __set(string $name, $value): void
    {
        if (property_exists($this, $name)) {
            $this->$name = $value;
        }
    }

    public function __isset(string $name): bool
    {
        return isset($this->$name);
    }

    const REQUIRED_ERROR_FORMAT_MESSAGE = '%s is required';
    const CHECK_EMPTY_RULE_FIELD = 'checkEmpty';
    protected function required($attribute, $options): bool
    {
        $value = $this->$attribute;
        if ($value === null) {
            $this->addError($attribute, sprintf(self::REQUIRED_ERROR_FORMAT_MESSAGE, $attribute));
            return false;
        }
        if (in_array(self::CHECK_EMPTY_RULE_FIELD, $options)) {
            if (empty($value)) {
                $this->addError($attribute, sprintf(self::REQUIRED_ERROR_FORMAT_MESSAGE, $attribute));
                return false;
            }
        }
        return true;
    }

    const NOT_AN_INTEGER_ERROR_FORMAT_MESSAGE = '%s should be integer';

    const STRICT_RULE_FIELD = 'strict';
    protected function integer($attribute, $options): bool
    {
        if (!is_numeric($this->$attribute)) {
            $this->addError($attribute, sprintf(self::NOT_AN_INTEGER_ERROR_FORMAT_MESSAGE, $attribute));
            return false;
        }
        if (in_array(self::STRICT_RULE_FIELD, $options)) {
            // do not accept '8' as correct integer field if rule is strict
            if (!is_integer($this->$attribute)) {
                $this->addError($attribute, sprintf(self::NOT_AN_INTEGER_ERROR_FORMAT_MESSAGE, $attribute));
                return false;
            }
        }
        $this->$attribute = (int)$this->$attribute;
        return true;
    }

    const NOT_A_STRING_ERROR_FORMAT_MESSAGE = '%s should be string';

    protected function string($attribute, $options): bool
    {
        // convert 8 to '8' if rule isn't strict
        if (!in_array(self::STRICT_RULE_FIELD, $options) && is_integer($this->$attribute)) {
            $this->$attribute = (string)$this->$attribute;
        }
        if (!is_string($this->$attribute)) {
            $this->addError($attribute, sprintf(self::NOT_A_STRING_ERROR_FORMAT_MESSAGE, $attribute));
            return false;
        }
        return true;
    }

}
