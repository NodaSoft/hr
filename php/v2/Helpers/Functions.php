<?php

namespace NW\WebService\References\Operations\Notification\Helpers;

function getResellerEmailFrom()
{
  // Реализация должна быть заменена на реальную логику
  return 'contractor@example.com';
}

function getEmailsByPermit($resellerId, $event)
{
  // Реализация должна быть заменена на реальную логику
  return ['someemeil@example.com', 'someemeil2@example.com'];
}

class NotificationEvents
{
  public const CHANGE_RETURN_STATUS = 'changeReturnStatus';
  public const NEW_RETURN_STATUS = 'newReturnStatus';
}

abstract class ReferencesOperation
{
  abstract public function doOperation(): array;

  protected function getRequest($pName)
  {
    // Реализация должна быть заменена на реальную логику
    return $_REQUEST[$pName] ?? null;
  }
}