<?php

declare(strict_types=1);


namespace NW\WebService\Request\DTO;

use NW\WebService\Notification\NotificationTypeEnum;

class RequestDTO
{

    public NotificationTypeEnum $notificationType;//required
    public ?PositionDTO $differences = null;
    public ?int $resellerId = null;
    public ?int $clientId = null;
    public ?int $creatorId = null;
    public ?int $expertId = null;
    public int $complaintId;//required
    public int $complaintNumber;//required
    public int $consumptionId;//required
    public int $consumptionNumber;//required
    public int $agreementNumber;//required
    public int $date;//required


}