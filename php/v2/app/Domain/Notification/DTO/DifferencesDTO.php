<?php

namespace app\Domain\Notification\DTO;

use Spatie\DataTransferObject\DataTransferObject;

class DifferencesDTO extends DataTransferObject
{
    /** @var int */
    public $notification_type;

    /** @var int */
    public $user_id;

    /** @var array */
    public $differences;
}