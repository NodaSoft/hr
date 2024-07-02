<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Dto;

/**
 * Class OperationResponseDto.
 *
 */
class OperationResponseDto
{
    /**
     * @var bool
     */
    private bool $notificationEmployeeByEmail;
    /**
     * @var bool
     */
    private bool $notificationClientByEmail;
    /**
     * @var array
     */
    private array $notificationClientBySms;

    /**
     * OperationResponseDto constructor.
     *
     * @param bool $notificationEmployeeByEmail
     * @param bool $notificationClientByEmail
     * @param array $notificationClientBySms
     */
    public function __construct(
        bool  $notificationEmployeeByEmail = false,
        bool  $notificationClientByEmail = false,
        array $notificationClientBySms = []
    )
    {
        $this->notificationEmployeeByEmail = $notificationEmployeeByEmail;
        $this->notificationClientByEmail = $notificationClientByEmail;

        $this->notificationClientBySms = $this->validateNotificationClientBySms($notificationClientBySms);
    }

    /**
     * Validate or set default value for NotificationClientBySms.
     *
     * @param array $notificationClientBySms
     * @return array
     */
    private function validateNotificationClientBySms(array $notificationClientBySms): array
    {
        if (is_bool($notificationClientBySms['isSent'] ?? null) &&
            ($notificationClientBySms['message'] ?? false)) {
            return $notificationClientBySms;
        }

        return [
            'isSent' => false,
            'message' => ''
        ];
    }

    /**
     * Set value for NotificationClientBySms isSent to true.
     *
     * @return void
     */
    public function setIsSent(): void
    {
        $this->notificationClientBySms['isSent'] = true;
    }

    /**
     * Set value for notificationEmployeeByEmail o true.
     *
     * @return void
     */
    public function setNotificationEmployeeByEmail(): void
    {
        $this->notificationEmployeeByEmail = true;
    }

    /**
     * Set value for NotificationClientBySms Message.
     *
     * @param string $message
     * @return void
     */
    public function setMessage(string $message): void
    {
        $this->notificationClientBySms['message'] = $message;
    }

    /**
     * Convert OperationResponseDto to array.
     *
     * @return array
     */
    public function toArray(): array
    {
        return [
            'notificationEmployeeByEmail' => $this->notificationEmployeeByEmail,
            'notificationClientByEmail' => $this->notificationClientByEmail,
            'notificationClientBySms' => $this->notificationClientBySms
        ];
    }
}