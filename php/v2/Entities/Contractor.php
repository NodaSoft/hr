<?php

namespace NW\WebService\References\Operations\Notification\Entities;

class Contractor
{
    const TYPE_CUSTOMER = 0;
    private string $id;
    private string $type;
    private string $name;
    private ?string $mobile;
    private ?string $email;
    private Seller $seller;

    public static function getById(int $resellerId): ?self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name . ' ' . $this->id;
    }


    public function getName(): string
    {
        return $this->name;
    }

    public function setName(string $name): self
    {
        $this->name = $name;

        return $this;
    }


    public function getId(): int
    {
        return $this->id;
    }


    public function setId(int $id): self
    {
        $this->id = $id;
        return $this;
    }


    public function setType(int $type): self
    {
        $this->type = $type;
        return $this;
    }


    public function getType(): int
    {
        return $this->type;
    }


    public function getMobile(): ?string
    {
        return $this->mobile;
    }


    public function setMobile(?string $mobile): self
    {
        $this->mobile = $mobile;
        return $this;
    }


    public function getEmail(): ?string
    {
        return $this->email;
    }


    public function setEmail(?string $email): self
    {
        $this->email = $email;
        return $this;
    }

    public function getSeller(): Seller
    {
        return $this->seller;
    }


    public function setSeller(Seller $seller): self
    {
        $this->seller = $seller;
        return $this;
    }

}