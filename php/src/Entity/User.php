<?php

namespace App\Entity;

use App\ORM;
use App\Repository\UserRepository;

#[ORM\Entity(repository: UserRepository::class, table: 'Users')]
class User
{
    #[ORM\ID]
    #[ORM\Column]
    private ?int $id = null;

    #[ORM\Column]
    private string $name;

    #[ORM\Column]
    private ?string $lastName;

    #[ORM\Column(type: ORM\ColumnType::INT, length: 3)]
    private int $age;

    #[ORM\Column]
    private ?string $_from;

    #[ORM\Column(ORM\ColumnType::JSON)]
    private $settings;

    public function getId(): ?int
    {
        return $this->id;
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

    public function getLastName(): string
    {
        return $this->lastName;
    }

    public function setLastName(?string $lastName): self
    {
        $this->lastName = $lastName;
        return $this;
    }

    public function getAge(): int
    {
        return $this->age;
    }

    public function setAge(string $age): self
    {
        $this->age = $age;
        return $this;
    }

    public function getFrom(): string
    {
        return $this->_from;
    }

    public function setFrom(string $from): self
    {
        $this->_from = $from;
        return $this;
    }

    public function getSettings(): ?array
    {
        return $this->settings;
    }

    public function setSettings(array $settings): self
    {
        $this->settings = $settings;
        return $this;
    }
}