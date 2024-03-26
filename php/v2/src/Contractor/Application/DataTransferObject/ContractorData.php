<?php

namespace Src\Contractor\Application\DataTransferObject;

use Src\Contractor\Domain\Entity\Contractor;

class ContractorData
{
    public int $id;
    public int $type;
    public string $name;
    public ?string $sellerId;
    public ?string $email;
    public ?string $mobile;

    public static function fromArray(array $data): self
    {
        $dto = new self();
        $dto->id = $data['id'];
        $dto->type = $data['type'];
        $dto->name = $data['name'];
        $dto->sellerId = $data['sellerId'];
        $dto->email = $data['email'] ?? null;
        $dto->mobile = $data['mobile'] ?? null;

        return $dto;
    }

    public static function fromEntity(Contractor $contractor): self
    {
        $dto = new self();
        $dto->id = $contractor->id;
        $dto->type = $contractor->type;
        $dto->name = $contractor->name;
        $dto->sellerId = $contractor->sellerId;
        $dto->email = $contractor->email;
        $dto->mobile = $contractor->mobile;

        return $dto;
    }

    public function toArray(): array
    {
        return [
            'id' => $this->id,
            'type' => $this->type,
            'name' => $this->name,
            'sellerId' => $this->sellerId,
            'email' => $this->email,
            'mobile' => $this->mobile,
        ];
    }

}