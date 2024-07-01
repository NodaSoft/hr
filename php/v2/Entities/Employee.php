<?php

namespace NW\WebService\References\Operations\Notification\Entities;

class Employee extends Contractor
{
  public static function getById(int $id): ?self
  {
    // Реализация должна быть заменена на реальную логику
    $employee       = new self();
    $employee->id   = $id;
    $employee->type = self::TYPE_CUSTOMER;
    $employee->name = "Employee $id";
    return $employee;
  }
}