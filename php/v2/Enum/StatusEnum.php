<?php

namespace NW\WebService\References\Operations\Notification;

enum StatusEnum: int
{
    case Completed = 0;
    case Pending = 1;
    case Rejected = 2;
}
