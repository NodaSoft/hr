<?php

use Src\Employee\Domain\Entity\Employee;

class EmployeeData
{
    public int $id;
    public int $type;
    public string $name;

    public static function fromArray(array $data): self
    {
        $dto = new self();
        $dto->id = $data['id'];
        $dto->type = $data['type'];
        $dto->name = $data['name'];

        return $dto;
    }

    public static function fromEntity(Employee $employee): self
    {
        $dto = new self();
        $dto->id = $employee->id;
        $dto->type = $employee->type;
        $dto->name = $employee->name;

        return $dto;
    }

    public function toArray(): array
    {
        return [
            'id' => $this->id,
            'type' => $this->type,
            'name' => $this->name,
        ];
    }
}