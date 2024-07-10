<?php

class TemplateDataIsEmptyException extends \Exception
{
    public function __construct(string $key)
    {
        parent::__construct("Template Data ({$key}) is empty!", 500);
    }
}