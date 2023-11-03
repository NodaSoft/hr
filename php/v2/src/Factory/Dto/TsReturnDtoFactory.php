<?php

namespace NodaSoft\Factory\Dto;

use NodaSoft\Dto\TsReturnDto;
use NodaSoft\ReferencesOperation\Params\TsReturnOperationParams;

class TsReturnDtoFactory
{
    public function makeTsReturnDto(TsReturnOperationParams $params): TsReturnDto
    {
        $dto = new TsReturnDto();
        foreach ($params->toArray() as $key => $value) {
            $setter = 'set' . $key;
            if (method_exists($dto, $setter)) {
                $dto->$setter($value);
            }
        }
        return $dto;
    }
}
