<?php

namespace Gateway;

class Condition
{
    private string $condition;
    private array $params;

    public function __construct(string $condition, array $params)
    {
        $this->condition = $condition;
        $this->params = $params;
    }

    /**
     * @return string
     */
    public function condition(): string
    {
        return $this->condition;
    }

    /**
     * @return array
     */
    public function params(): array
    {
        return $this->params;
    }
}