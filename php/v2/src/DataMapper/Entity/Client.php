<?php

namespace NodaSoft\DataMapper\Entity;

class Client implements Entity
{
    /** @var int */
    private $id;

    /** @var string */
    private $name;

    /** @var Reseller */
    private $reseller;

    /** @var string */
    private $email;

    /** @var int */
    private $cellphoneNumber;

    /** @var bool */
    private $isCustomer;

    public function setId(int $id): void
    {
        $this->id = $id;
    }

    public function getId(): int
    {
        return $this->id;
    }

    public function setName(string $name): void
    {
        $this->name = $name;
    }

    public function getName(): string
    {
        return $this->name;
    }

    public function getFullName(): string
    {
        return $this->getName() . ' ' . $this->getId();
    }

    public function getReseller(): Reseller
    {
        return $this->reseller;
    }

    public function setReseller(Reseller $reseller): void
    {
        $this->reseller = $reseller;
    }

    public function hasReseller(Reseller $reseller): bool
    {
        return $this->reseller->getId() === $reseller->getId();
    }

    public function hasEmail(): bool
    {
        return is_null($this->getEmail());
    }

    public function getEmail(): ?string
    {
        return $this->email ?? null;
    }

    public function setEmail(string $email): void
    {
        $this->email = $email;
    }

    public function hasCellphoneNumber(): bool
    {
        return is_null($this->getCellphoneNumber());
    }

    public function getCellphoneNumber(): ?int
    {
        return $this->cellphoneNumber ?? null;
    }

    public function setCellphoneNumber(int $cellphoneNumber): void
    {
        $this->cellphoneNumber = $cellphoneNumber;
    }

    public function isCustomer(): bool
    {
        return $this->isCustomer;
    }

    public function setIsCustomer(bool $isCustomer): void
    {
        $this->isCustomer = $isCustomer;
    }
}
