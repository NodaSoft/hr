<?php

declare(strict_types=1);


namespace NW\WebService\Request;

use NW\WebService\Position\PositionStatusEnum;
use NW\WebService\Request\DTO\PositionDTO;
use NW\WebService\Request\DTO\RequestDTO;


//конечно можно описать интерфейсы, класс request с валидацией, но обошелся малой кровью

class RequestMapper
{

    public static function fromPost(): RequestDTO
    {
        return self::mapFromArray($_POST['data']);
    }

    public static function fromGet(): RequestDTO
    {
        return self::mapFromArray($_GET['data']);
    }

    private static function mapFromArray(array $data): RequestDTO
    {
        $dto = new RequestDTO();

        $dto->resellerId = $data['resellerId'];
        $dto->notificationType = $data['notificationType']; //предположим тут NotificationTypeEnum
        $dto->clientId = $data['clientId'];
        $dto->creatorId = $data['creatorId'];
        $dto->expertId = $data['expertId'];

        if ($data['differences'] ?? 0) {//уже на входе должны быть определены enum-ы
            $dto->differences = new PositionDTO(
                from: PositionStatusEnum::COMPLETED,
                to: PositionStatusEnum::REJECTED
            );
        }

        $dto->complaintId = $data['complaintId'];
        $dto->complaintNumber = $data['complaintNumber'];
        $dto->consumptionId = $data['consumptionId'];
        $dto->consumptionNumber = $data['consumptionNumber'];
        $dto->agreementNumber = $data['agreementNumber'];
        $dto->date = $data['date'];


        return $dto;
    }
}