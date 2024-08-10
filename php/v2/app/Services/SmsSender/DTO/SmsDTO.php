<?php

namespace app\Services\SmsSender\DTO;

use Spatie\DataTransferObject\DataTransferObject;
class SmsDTO extends DataTransferObject
{
    /** @var array */
    public $data;

    /** @var int */
    public $user_id;

    /** @var string */
    public $status;

    /** @var int|null */
    public $client_id;

    /** @var int|null */
    public $to;
}