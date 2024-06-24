<?php

namespace NW\WebService\References\Operations\Notification;

readonly final class EnumRule
{
    public function __construct(
        private string $type
    ) {
    }

    public function __invoke($value)
    {
        if (!enum_exists($this->type) || !method_exists($this->type, '__toString')) {
            return false;
        }

        return is_null($this->type::tryFrom($value)) ? false: $value;
    }
}
