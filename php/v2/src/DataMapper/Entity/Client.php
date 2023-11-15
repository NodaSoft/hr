<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\Messenger\Recipient;
use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class Client implements Entity, Recipient
{
    use EntityTrait\MessageRecipientEntity;

    /** @var Reseller */
    private $reseller; //todo: should I rename reseller to seller?

    /** @var int */
    private $cellphoneNumber;

    /** @var bool */
    private $isCustomer;

    public function __construct(
        int $id = null,
        string $name = null,
        string $email = null,
        int $cellphone = null,
        bool $isCustomer = null,
        Reseller $reseller = null
    ) {
        if ($id) $this->setId($id);
        if ($name) $this->setName($name);
        if ($email) $this->setEmail($email);
        if ($cellphone) $this->setCellphone($cellphone);
        if (is_bool($isCustomer)) $this->setIsCustomer($isCustomer);
        if ($reseller) $this->setReseller($reseller);
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
        return isset($this->reseller)
            && $this->reseller->getId() === $reseller->getId();
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
