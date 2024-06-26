<?php

namespace Nodasoft\Testapp\Entities\Client;

use Nodasoft\Testapp\Entities\Seller\Seller;
use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Interfaces\ContactorInterface;

readonly class Client implements ContactorInterface
{
    public function __construct(
        private int           $id,
        private string        $name,
        private int           $sellerId,
        private ?string       $email,
        private ?string       $mobile,
        private ContactorType $type
    )
    {
    }


    public function getId(): int
    {
        return $this->id;
    }

    public function getSellerId(): int
    {
        return $this->sellerId;
    }

    public function getMobile(): ?string
    {
        return $this->mobile;
    }

    public function getName(): string
    {
        return $this->name;
    }

    public function getEmail(): ?string
    {
        return $this->email;
    }

    public function getType(): ContactorType
    {
        return $this->type;
    }

    public function getFullName(): string
    {
        return $this->name . ' lastname';
    }

    public function seller(): Seller
    {
        return new Seller(
            $this->sellerId,
            'has one seller',
            'has one seller',
            ContactorType::TYPE_SELLER
        );
    }

    public function isCustomer(): bool
    {
        return $this->type === ContactorType::TYPE_CUSTOMER;
    }

    public function HasSeller(int $sellerId): bool
    {
        return $this->seller()->getId() === $sellerId;
    }
}