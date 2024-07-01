<?php

namespace NW\WebService\References\Operations\Notification\Validators;

use NW\WebService\References\Operations\Notification\Exceptions\ValidationException;

class DataValidator
{
  public function validateData(array $data): void
  {
    if (empty((int) $data['resellerId'])) {
      throw new ValidationException('Empty resellerId', 400);
    }

    if (empty((int) $data['notificationType'])) {
      throw new ValidationException('Empty notificationType', 400);
    }

    // Больше правил...
  }

  public function validateTemplateData(array $templateData): void
  {
    foreach ($templateData as $key => $value) {
      if (empty($value)) {
        throw new ValidationException("Template Data ({$key}) is empty!", 500);
      }
    }
  }
}