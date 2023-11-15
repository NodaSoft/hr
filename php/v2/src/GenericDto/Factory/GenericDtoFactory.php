<?php

namespace NodaSoft\GenericDto\Factory;

use NodaSoft\GenericDto\Dto\Dto;
use NodaSoft\Operation\Params\Params;

class GenericDtoFactory
{
    /**
     * @template ParticularDTO of Dto
     * @param ParticularDTO $dto
     * @param Params $params
     * @return ParticularDTO
     */
    public function fillDtoParams(
        Dto                       $dto,
        Params $params
    ): Dto {
        foreach ($params->toArray() as $key => $value) {
            $setter = 'set' . $key;
            if (method_exists($dto, $setter)) {
                $dto->$setter($value);
            }
        }
        return $dto;
    }

    /**
     * @template ParticularDTO of Dto
     * @param ParticularDTO $dto
     * @param Params $params
     * @return ParticularDTO
     */
    public function fillDtoArray(
        Dto      $dto,
        iterable $params
    ): Dto {
        foreach ($params as $key => $value) {
            $setter = 'set' . $key;
            if (method_exists($dto, $setter)) {
                $dto->$setter($value);
            }
        }
        return $dto;
    }
}
