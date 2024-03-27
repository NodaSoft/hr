<?php

namespace App\Http\Request;

use App\Enum\Notification;
use App\Enum\Status;
use DateTimeInterface;
use phpDocumentor\Reflection\Types\This;
use Symfony\Component\Validator\Constraints as Assert;

use Symfony\Component\Validator\Exception\ValidationFailedException;
use Symfony\Component\Validator\Validation;
use Symfony\Component\HttpFoundation\Request;

class StatusChangedRequest extends BaseRequest
{
    public const NEW_STATUS = 'new_status';

    protected function getRules(): array
    {
        $dataRules = $this->getBaseDataRules();
        $dataRules[self::NEW_STATUS] = [
            new Assert\NotBlank(),
            new Assert\Positive(),
            new Assert\Choice([
                Status::COMPLETED->value,
                Status::PENDING->value,
                Status::REJECTED->value,
            ])
        ];

        return [
            self::DATA => $dataRules
        ];
    }
}
