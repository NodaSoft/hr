<?php

namespace src\Operation\Infrastructure\Domain\Enum;

enum ContractorType: int
{
    const TYPE_CUSTOMER = 0;
    const TYPE_SELLER = 1;
    const TYPE_EMPLOYEE = 2;
}
