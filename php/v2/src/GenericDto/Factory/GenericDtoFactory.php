<?php

namespace NodaSoft\GenericDto\Factory;

use NodaSoft\GenericDto\Dto\DtoInterface;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;

class GenericDtoFactory
{
    /**
     * @template ParticularDTO of DtoInterface
     * @param ParticularDTO $dto
     * @param ReferencesOperationParams $params
     * @return ParticularDTO
     */
    public function fillDtoParams(
        DtoInterface $dto,
        ReferencesOperationParams $params
    ): DtoInterface {
        foreach ($params->toArray() as $key => $value) {
            $setter = 'set' . $key;
            if (method_exists($dto, $setter)) {
                $dto->$setter($value);
            }
        }
        return $dto;
    }

    /**
     * @template ParticularDTO of DtoInterface
     * @param ParticularDTO $dto
     * @param ReferencesOperationParams $params
     * @return ParticularDTO
     */
    public function fillDtoArray(
        DtoInterface $dto,
        iterable $params
    ): DtoInterface {
        foreach ($params as $key => $value) {
            $setter = 'set' . $key;
            if (method_exists($dto, $setter)) {
                $dto->$setter($value);
            }
        }
        return $dto;
    }
}
