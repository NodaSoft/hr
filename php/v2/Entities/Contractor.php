<?php

namespace NW\WebService\References\Operations\Notification\Entities;

class Contractor
{
  public const TYPE_CUSTOMER = 0;

  public int $id;
  public int $type;
  public string $name;
  public ?string $email;
  public ?string $mobile;

  public static function getById(int $id): ?self
  {
    // Реализация должна быть заменена на реальную логику
    $contractor       = new self();
    $contractor->id   = $id;
    $contractor->type = self::TYPE_CUSTOMER;
    $contractor->name = "Contractor $id";
    return $contractor;
  }

  public function getFullName(): string
  {
    return $this->name . ' ' . $this->id;
  }
}