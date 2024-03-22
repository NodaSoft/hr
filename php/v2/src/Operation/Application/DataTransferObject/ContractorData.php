<?php

namespace Src\Operation\Application\DataTransferObject;

class ContractorData
{
    public int $id;
    public int $type;
    public string $name;

    public string $fullName;
    public ?int $sellerId;

    /**
     * @param array $data
     * @return self
     */
    public static function fromArray(array $data): self
    {
        $dto = new self();
        $dto->id = $data['id'];
        $dto->type = $data['type'];
        $dto->name = $data['name'];
        $dto->fullName = $data['name'] . ' ' . $data['id'];
        $dto->sellerId = $data['sellerId'] ?? null;

        return $dto;
    }

    /**
     * @return array
     */
    public function toArray(): array
    {
        return [
            'id' => $this->id,
            'type' => $this->type,
            'name' => $this->name,
            'fullName' => $this->fullName,
            'sellerId' => $this->sellerId,
        ];
    }

}