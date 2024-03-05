<?php

namespace NW\WebService\Position;

enum PositionStatusEnum: int
{

    case COMPLETED = 1;
    case PENDING = 2;
    case REJECTED = 3;
}
