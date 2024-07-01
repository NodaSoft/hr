<?php

namespace NW\WebService\References\Operations\Notification\Entities;

class Seller extends Contractor
{
  public static function getById(int $id): ?self
  {
    // Реализация должна быть заменена на реальную логику
    $seller       = new self();
    $seller->id   = $id;
    $seller->type = self::TYPE_CUSTOMER;
    $seller->name = "Seller $id";
    return $seller;
  }
}