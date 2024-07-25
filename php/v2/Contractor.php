<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */
class Contractor
{
    const TYPE_CUSTOMER = 0;
    protected int $id;
    protected ?int $type = null;
    protected ?string $name = null;
    protected ?string $email = null;
    protected ?string $mobile = null;

    public function __construct(int $id)
    {
        $this->id = $id;
    }

    public static function getById(int $resellerId): self
    {
        return new self($resellerId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return sprintf('%s %s', $this->name, $this->id);
    }

    public function getType(): int
    {
        return $this->type;
    }

    public function getId(): int
    {
        return $this->id;
    }

    public function getName(): ?string
    {
        return $this->name;
    }

    public function getEmail(): ?string
    {
        return $this->email;
    }

    public function getMobile(): ?string
    {
        return $this->mobile;
    }
}