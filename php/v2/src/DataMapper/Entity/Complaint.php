<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;

class Complaint implements Entity
{
    use EntityTrait\Entity;

    /** @var Employee */
    private $creator;

    /** @var Client */
    private $client;

    /** @var Employee */
    private $expert;

    /** @var Reseller */
    private $reseller;

    /** @var ComplaintStatus */
    private $status;

    /** @var ?ComplaintStatus */
    private $previousStatus = null;

    /** @var string */
    private $number;

    public function __construct(
        int $id = null,
        string $name = null,
        Employee $creator = null,
        Client $client = null,
        Employee $expert = null,
        Reseller $reseller = null,
        ComplaintStatus $status = null,
        ComplaintStatus $previousStatus = null,
        string $number = null
    ) {
        if ($id) $this->setId($id);
        if ($name) $this->setName($name);
        if ($creator) $this->setCreator($creator);
        if ($client) $this->setClient($client);
        if ($expert) $this->setExpert($expert);
        if ($reseller) $this->setReseller($reseller);
        if ($status) $this->setStatus($status);
        if ($previousStatus) $this->setPreviousStatus($previousStatus);
        if ($number) $this->setNumber($number);
    }

    public function setCreator(Employee $creator): void
    {
        $this->creator = $creator;
    }

    public function setClient(Client $client): void
    {
        $this->client = $client;
    }

    public function setStatus(ComplaintStatus $status): void
    {
        $this->status = $status;
    }

    public function setExpert(Employee $expert): void
    {
        $this->expert = $expert;
    }

    public function setReseller(Reseller $reseller): void
    {
        $this->reseller = $reseller;
    }

    public function setPreviousStatus(ComplaintStatus $previousStatus): void
    {
        $this->previousStatus = $previousStatus;
    }

    public function getCreator(): Employee
    {
        return $this->creator;
    }

    public function getClient(): Client
    {
        return $this->client;
    }

    public function getExpert(): Employee
    {
        return $this->expert;
    }

    public function getReseller(): Reseller
    {
        return $this->reseller;
    }

    public function getStatus(): ComplaintStatus
    {
        return $this->status;
    }

    public function getPreviousStatus(): ?ComplaintStatus
    {
        return $this->previousStatus;
    }

    public function getNumber(): string
    {
        return $this->number;
    }

    public function setNumber(string $number): void
    {
        $this->number = $number;
    }
}
