<?php

namespace NW\WebService\References\Operations\Notification\Request;

use Exception;

readonly abstract class RequestAbstract
{
    abstract protected static function rules(): array;

    /**
     * @throws Exception
     */
    public static function fromArray($data): static
    {
        if (!is_array($data)) {
            throw new Exception('Data must be an array');
        }

        return new static(...self::validate($data));
    }

    /**
     * @throws Exception
     */
    private static function validate(array $data): array
    {
        $rules = static::rules();

        $params = array_map(fn($k, $v) => [$k => $rules[$k]['params']], array_keys($rules), $rules);
        $dataValidated = filter_var_array($data, $params);

        if (!$dataValidated) {
            throw new Exception('Wrong data provided');
        }

        $none = array_filter(
            $dataValidated,
            fn($v, $k) => $v === null && !empty($rules[$k]['required']), ARRAY_FILTER_USE_BOTH
        );

        if ($none) {
            throw new Exception(sprintf('Column(s) %s empty!', implode(', ', array_keys($none))), 400);
        }

        $notValidated = array_filter(
            $dataValidated,
            fn($v, $k) =>
                !empty($rules[$k]['flags']) && $rules[$k]['flags'] === FILTER_REQUIRE_ARRAY ? in_array(false, $v) : $v === false,
            ARRAY_FILTER_USE_BOTH
        );

        if ($notValidated) {
            throw new Exception(sprintf('Column(s) %s empty!', implode(', ', array_keys($notValidated))), 400);
        }

        return $dataValidated;
    }
}