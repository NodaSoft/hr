<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification\Dto;

/**
 * Class EmailMessageDto.
 *
 */
class EmailMessageDto
{
    /**
     * @var string
     */
    private string $emailFrom;
    /**
     * @var string
     */
    private string $emailTo;
    /**
     * @var string
     */
    private string $subject;
    /**
     * @var string
     */
    private string $message;

    public function __construct(
        string $emailFrom,
        string $emailTo,
        string $subject,
        string $message
    )
    {
        $this->emailFrom = $emailFrom;
        $this->emailTo = $emailTo;
        $this->subject = $subject;
        $this->message = $message;
    }

    public function toArray(): array
    {
        return [
            'emailFrom' => $this->emailFrom,
            'emailTo' => $this->emailTo,
            'subject' => $this->subject,
            'message' => $this->message
        ];
    }
}