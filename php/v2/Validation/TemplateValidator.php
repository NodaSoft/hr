<?php

namespace NW\WebService\References\Operations\Notification\Validation;

/**
 * TemplateValidator class
 */
class TemplateValidator implements ValidatorInterface
{

    public function validate(array $data, array &$result = []): bool {
        foreach ($data as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        return true;
    }
}
