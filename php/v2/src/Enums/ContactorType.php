<?php

namespace Nodasoft\Testapp\Enums;

enum ContactorType: int
{
    case TYPE_CUSTOMER = 1;
    case TYPE_EMPLOYEE = 2;
    case TYPE_SELLER = 3;
}