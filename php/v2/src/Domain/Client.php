<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Domain;

use Exception;
use NW\WebService\References\Operations\Notification\Enum\ContractorType;
use Symfony\Component\HttpFoundation\Response;

final class Client extends Contractor
{
    /**
     * @var non-empty-string
     */
    public string $mobile;

    public static function find($contractorId): self
    {
        $contractor = parent::find($contractorId);

        /** @phpstan-ignore-next-line */
        if ($contractor->type != ContractorType::TYPE_CUSTOMER) {
            throw new Exception('Client not found!', Response::HTTP_BAD_REQUEST);
        }

        return new self($contractorId);
    }
}
