<?php

declare(strict_types=1);

namespace ResultOperation\Entity;

use pseudovendor\BaseEntity;
use ResultOperation\Enum\ContractorType;

class Contractor extends BaseEntity
{
    private const TYPES = [
        0 => ContractorType::DEFAULT,
        1 => ContractorType::CUSTOMER,
    ];

    /**
     * @inheritdoc
     */
    protected array $fields = [
        'id',
        'name',
        'type',
        'seller',
        'email',
        'has_mobile'
    ];

    /**
     * @return string
     */
    public function getFullName(): string
    {
        return trim(
            sprintf(
                '%s %d',
                $this->getAttribute('name') ?? '',
                $this->getId()
            )
        );
    }

    /**
     * @return ?int
     */
    public function getSellerId(): ?int
    {
        return $this->getAttribute('seller');
    }

    /**
     * @return ?ContractorType
     */
    public function getType(): ?ContractorType
    {
        return $this->getAttribute('type');
    }

    /**
     * @return ?string
     */
    public function getEmail(): ?string
    {
        return $this->getAttribute('email');
    }

    /**
     * @return bool
     */
    public function hasMobile(): bool
    {
        return (bool) $this->getAttribute('has_mobile');
    }

    /**
     * "Волшебный" метод фреймворка, который значение из бд подменит на enum
     *
     * @param int $type
     * @return self
     */
    protected function setTypeAttribute(int $type): self
    {
        return $this->setAttribute('type', self::TYPES[$type] ?? ContractorType::DEFAULT);
    }
}
