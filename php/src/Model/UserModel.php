<?php

declare(strict_types=1);

namespace App\Model;

class UserModel
{
    private string $login;
    private string $name;
    private string $lastName;
    private int $age;
    private ?string $from = null;
    private ?array $settings = null;

    /**
     * @return string
     */
    public function getLogin(): string
    {
        return $this->login;
    }

    /**
     * @param string $login
     * @return UserModel
     */
    public function setLogin(string $login): UserModel
    {
        $this->login = $login;
        return $this;
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }

    /**
     * @param string $name
     * @return UserModel
     */
    public function setName(string $name): UserModel
    {
        $this->name = $name;
        return $this;
    }

    /**
     * @return string
     */
    public function getLastName(): string
    {
        return $this->lastName;
    }

    /**
     * @param string $lastName
     * @return UserModel
     */
    public function setLastName(string $lastName): UserModel
    {
        $this->lastName = $lastName;
        return $this;
    }

    /**
     * @return int
     */
    public function getAge(): int
    {
        return $this->age;
    }

    /**
     * @param int $age
     * @return UserModel
     */
    public function setAge(int $age): UserModel
    {
        $this->age = $age;
        return $this;
    }

    /**
     * @return string|null
     */
    public function getFrom(): ?string
    {
        return $this->from;
    }

    /**
     * @param string|null $from
     * @return UserModel
     */
    public function setFrom(?string $from): UserModel
    {
        $this->from = $from;
        return $this;
    }

    /**
     * @return array|null
     */
    public function getSettings(): ?array
    {
        return $this->settings;
    }

    /**
     * @param array|null $settings
     * @return UserModel
     */
    public function setSettings(?array $settings): UserModel
    {
        $this->settings = $settings;
        return $this;
    }
}
