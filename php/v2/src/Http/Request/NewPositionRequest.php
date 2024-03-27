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

class NewPositionRequest extends BaseRequest
{
    protected function getRules(): array
    {
        $dataRules = $this->getBaseDataRules();

        return [
            self::DATA => $dataRules
        ];
    }
}
