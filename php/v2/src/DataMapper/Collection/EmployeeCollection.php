<?php

namespace NodaSoft\DataMapper\Collection;

use NodaSoft\DataMapper\Entity\Employee;

/**
 * @implements \Iterator<int, Employee>
 */
class EmployeeCollection implements \Iterator
{
    /** @var Employee[] */
    private $collection = [];

    /** @var int */
    private $pointer = 0;

    /**
     * @param Employee[] $employees
     * @throws \Exception only Employee[] allowed
     */
    public function __construct(array $employees = [])
    {
        foreach ($employees as $employee) {
            if (! $employee instanceof Employee) {
                throw new \Exception('Only ' . Employee::class . ' is allowed for the constructor.');
            }
            $this->collection[] = $employee;
        }
    }

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

    public function count(): int
    {
        return count($this->collection);
    }

    public function getEmployee(int $id): ?Employee
    {
        foreach ($this->collection as $employee) {
            if ($employee->getId() === $id) {
                return $employee;
            }
        }

        return null;
    }
}
