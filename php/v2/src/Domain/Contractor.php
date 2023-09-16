<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Domain;

use Exception;
use NW\WebService\References\Operations\Notification\Enum\ContractorType;
use Symfony\Component\HttpFoundation\Response;

class Contractor
{
    /**
     * @var positive-int
     */
    public int $id;

    public ContractorType $type;

    public string $name;

    public Seller $seller;

    /**
     * @param positive-int $contractorId
     */
    public function __construct(int $contractorId)
    {
        $this->id = $contractorId;
    }

    /**
     * @param positive-int $contractorId
     * @throws Exception
     */
    public static function find(int $contractorId): self
    {
        $contractor = new self($contractorId); // fakes the getById method

        if ($contractor == null) {
            throw new Exception('Contractor not found!', Response::HTTP_BAD_REQUEST);
        }

        return $contractor;
    }

    /**
     * @return non-empty-string
     */
    public function getFullName(): string
    {
        if ($this->name == '') {
            return (string) $this->id;
        }
        return $this->name . ' ' . $this->id;
    }
}
