<?php
namespace NW\WebService\References\Operations\Notification\Objects;
class NotificationClientBySms extends PlainOldObject
{
    public function __construct(
        public bool   $isSend = false,
        public string $message = ''
    )
    {}
}