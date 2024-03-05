<?php

declare(strict_types=1);

namespace NW\WebService\Employees;

class Employee
{
    private static array $storageEmployees = [];

    public function __construct(
        private readonly int $id,
        private readonly EmployeeTypeEnum $type,
        private readonly string $name,
        private readonly string $surname,
        private readonly string $email,
        private readonly string $phone
    ) {
    }

    public static function getById(EmployeeTypeEnum $type, int $id): ?Employee
    {
        return static::$storageEmployees[$type->value][$id] ?? null;
    }

    public static function setEmployee(Employee $employee): void
    {
        static::$storageEmployees[$employee->type->value][$employee->id] = $employee;
    }

    public function getFullName(): string
    {
        return $this->name.' '.$this->surname;
    }

    public function getType(): EmployeeTypeEnum
    {
        return $this->type;
    }

    public function getEmail(): string
    {
        return $this->email;
    }

    public function getPhone(): string
    {
        return $this->phone;
    }


    public function getId(): int
    {
        return $this->id;
    }


}
