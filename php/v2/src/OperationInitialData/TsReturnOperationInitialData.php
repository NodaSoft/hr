<?php

namespace NodaSoft\OperationInitialData;

use NodaSoft\Dto\TsReturnDto;

class TsReturnOperationInitialData implements OperationInitialData
{
    /** @var TsReturnDto */
    private $messageTemplate;

    /** @var int */
    private $resellerId;

    /** @var int */
    private $notificationType;

    public function getMessageTemplate(): TsReturnDto
    {
        return $this->messageTemplate;
    }

    public function setMessageTemplate(TsReturnDto $messageTemplate): void
    {
        $this->messageTemplate = $messageTemplate;
    }

    public function getResellerId(): int
    {
        return $this->resellerId;
    }

    public function setResellerId(int $resellerId): void
    {
        $this->resellerId = $resellerId;
    }

    public function getNotificationType(): int
    {
        return $this->notificationType;
    }

    public function setNotificationType(int $notificationType): void
    {
        $this->notificationType = $notificationType;
    }
}
