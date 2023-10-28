<?php

namespace NodaSoft\OperationInitialData;

use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\Dto\TsReturnDto;

class TsReturnOperationInitialData implements OperationInitialData
{
    /** @var TsReturnDto */
    private $messageTemplate;

    /** @var Reseller */
    private $reseller;

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

    public function getReseller(): Reseller
    {
        return $this->reseller;
    }

    public function setReseller(Reseller $reseller): void
    {
        $this->reseller = $reseller;
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
