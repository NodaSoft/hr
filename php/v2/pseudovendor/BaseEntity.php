<?php

declare(strict_types=1);

namespace pseudovendor;

/**
 * Условный библиотечный класс для работы с сущностью
 */
class BaseEntity
{
    /**
     * @var array
     */
    protected array $fields = [];

    /**
     * @param string $name
     * @return mixed
     */
    public function getAttribute(string $name): mixed
    {
        return 'value'; // stub
    }

    /**
     * @param string $param
     * @param mixed $value
     * @return self
     */
    protected function setAttribute(string $param, mixed $value): self
    {
        $this->fields[$param] = $value;

        return $this;
    }

    /**
     * @return int
     */
    public function getId(): int
    {
        return (int) $this->getAttribute('id'); // stub
    }
}
