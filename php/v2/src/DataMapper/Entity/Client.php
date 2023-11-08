<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\EmailEntity;
use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class Client implements Entity, EmailEntity
{
    use EntityTrait\EmailEntity;

    /** @var Reseller */
    private $reseller; //todo: should I rename reseller to seller?

    /** @var int */
    private $cellphoneNumber;

    /** @var bool */
    private $isCustomer;

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
        return isset($this->reseller)
            && $this->reseller->getId() === $reseller->getId();
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
