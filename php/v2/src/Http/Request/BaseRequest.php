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

class BaseRequest
{
    public const ID = 'id';
    public const DATA = 'data';
    public const NUMBER = 'number';
    public const RESELLER_ID = 'reseller_id';
    public const CLIENT_ID = 'client_id';
    public const CREATOR_ID = 'creator_id';
    public const EXPERT_ID = 'expert_id';
    public const COMPLAINT = 'complaint';
    public const CONSUMPTION = 'consumption';
    public const AGREEMENT_NUMBER = 'agreement_number';
    public const DATE = 'date';

    protected Request $request;

    public function __construct()
    {
        $this->request = Request::createFromGlobals();
        $this->validateRequest();
    }

    protected function getBaseDataRules(): array
    {
        return [
            self::RESELLER_ID => $this->positiveNotBlankRules(),
            self::CLIENT_ID => $this->positiveNotBlankRules(),
            self::CREATOR_ID => $this->positiveNotBlankRules(),
            self::CONSUMPTION => $this->idNumberRules(),
            self::EXPERT_ID => $this->positiveNotBlankRules(),
            self::COMPLAINT => $this->idNumberRules(),
            self::DATE => [new Assert\DateTime(DateTimeInterface::ATOM)],
            self::AGREEMENT_NUMBER => $this->positiveNotBlankRules(),
        ];
    }


    protected function validateRequest(): void
    {
        $validator = Validation::createValidator();
        $violations = $validator->validate($this->request->request->all(), $this->getRules());
        if ($violations->count() > 0) {
            throw new ValidationFailedException('request', $violations);
        }
    }

    public function getField(string $key, mixed $default = null)
    {
        return $this->request->get($key, $default);
    }


    public function getDataField(string $key, mixed $default = null)
    {
        $data = $this->request->get(self::DATA, []);

        return $data[$key] ?? $default;
    }

    protected function idNumberRules(): array
    {
        return [new Assert\Collection([
            self::ID => $this->positiveNotBlankRules(),
            self::NUMBER => $this->positiveNotBlankRules(),
        ])];
    }

    protected function positiveNotBlankRules(): array
    {
        return [new Assert\NotBlank(), new Assert\Positive()];
    }
}
