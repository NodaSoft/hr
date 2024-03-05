<?php

declare(strict_types=1);


namespace NW\WebService\Request\DTO;

use NW\WebService\Position\PositionStatusEnum;

class PositionDTO
{

    public function __construct(
        public PositionStatusEnum $from,
        public PositionStatusEnum $to
    ) {
    }
}