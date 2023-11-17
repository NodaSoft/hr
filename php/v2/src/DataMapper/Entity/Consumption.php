<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class Consumption implements Entity
{
    use EntityTrait\Entity;

    /** @var string */
    private $number;

    /** @var string */
    private $agreementNumber;

    public function __construct(
        int $id = null,
        string $name = null
    ) {
        if ($id) $this->setId($id);
        if ($name) $this->setName($name);
    }

    public function getNumber(): string
    {
        return $this->number;
    }

    public function setNumber(string $number): void
    {
        $this->number = $number;
    }

    public function getAgreementNumber(): string
    {
        return $this->agreementNumber;
    }

    public function setAgreementNumber(string $agreementNumber): void
    {
        $this->agreementNumber = $agreementNumber;
    }
}
