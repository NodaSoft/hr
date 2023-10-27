<?php

namespace NodaSoft\Factory\Dto;

use NodaSoft\Dto\TsReturnDto;

class TsReturnDtoFactory
{
    public function makeTsReturnDto(array $requestData): TsReturnDto
    {
        $dto = new TsReturnDto();
        foreach ($requestData as $key => $value) {
            $setter = 'set' . $key;
            if (method_exists($dto, $setter)) {
                $dto->$setter($value);
            }
        }
        return $dto;
    }
}
