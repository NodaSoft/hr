<?php
namespace NW\WebService\References\Operations\Notification\Objects;
class ReturnOperationResult extends PlainOldObject
{
    /**
     * @param bool $notificationEmployeeByEmail
     * @param bool $notificationClientByEmail
     * @param NotificationClientBySms|null $notificationClientBySms
     */
    public function __construct(
        public bool $notificationEmployeeByEmail = false,
        public bool $notificationClientByEmail = false,
        public ?NotificationClientBySms $notificationClientBySms = null
    )
    {
        $this->notificationClientBySms = $notificationClientBySms ?? new NotificationClientBySms();
    }
}