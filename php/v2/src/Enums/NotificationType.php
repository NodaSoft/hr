<?php

namespace Nodasoft\Testapp\Enums;

enum NotificationType: int
{
    case TYPE_NEW = 1;
    case TYPE_CHANGE = 2;

    public function isNew(): bool
    {
        return $this === NotificationType::TYPE_NEW;
    }

    public function isChange(): bool
    {
        return $this === NotificationType::TYPE_CHANGE;
    }
}