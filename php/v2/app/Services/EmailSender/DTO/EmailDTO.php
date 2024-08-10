<?php

namespace app\Services\EmailSender\DTO;

use Spatie\DataTransferObject\DataTransferObject;
class EmailDTO extends DataTransferObject
{
    /** @var string */
    public $email_from;

    /** @var string */
    public $email_to;

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