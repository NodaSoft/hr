<?php

namespace NodaSoft\DataMapper\Collection;

use NodaSoft\DataMapper\Entity\Employee;

class EmployeeCollection implements \Iterator
{
    /** @var Employee[] */
    private $collection = [];

    /** @var int */
    private $pointer = 0;

    public function add(Employee $employee): void
    {
        $this->collection[] = $employee;
    }

    public function current(): Employee
    {
        return $this->collection[$this->pointer];
    }

    public function next(): void
    {
        ++ $this->pointer;
    }

    public function key(): int
    {
        return $this->pointer;
    }

    public function valid(): bool
    {
        return isset($this->collection[$this->pointer]);
    }

    public function rewind(): void
    {
        $this->pointer = 0;
    }
}
