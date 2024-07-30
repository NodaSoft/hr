<?php

namespace NW\WebService\Employees;

enum EmployeeTypeEnum: int
{
    case RESELLER = 1;
    case CONTRACTOR = 2;
    case CREATOR = 3;
    case EXPERT = 4;
}
