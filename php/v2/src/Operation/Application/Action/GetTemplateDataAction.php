<?php

namespace Src\Operation\Application\Action;

use Src\Operation\Application\DataTransferObject\ContractorData;
use Src\Operation\Application\DataTransferObject\EmployeeData;
use Src\Operation\Application\DataTransferObject\OperationData;
use Src\Operation\Application\DataTransferObject\SellerData;
use src\Operation\Domain\Enum\NotificationType;
use src\Operation\Domain\Enum\PositionStatus;

class GetTemplateDataAction
{
    public function execute(
        OperationData  $operationData,
        SellerData     $reseller,
        EmployeeData   $creator,
        EmployeeData   $expert,
        ContractorData $client
    ): array
    {
        $differences = '';
        if ($operationData->notificationType === NotificationType::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $reseller->id);
        }
        if ($operationData->notificationType === NotificationType::TYPE_CHANGE && !empty($operationData->differences)) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => PositionStatus::STATUSES[$operationData->differences['from']],
                'TO' => PositionStatus::STATUSES[$operationData->differences['to']],
            ], $reseller
                ->id);
        }

        //   Здесь можно разместить дополнительно валидацию

        return [
            'COMPLAINT_ID' => $operationData->complaintId,
            'COMPLAINT_NUMBER' => $operationData->complaintNumber,
            'CREATOR_ID' => $operationData->creatorId,
            'CREATOR_NAME' => $creator->fullName,
            'EXPERT_ID' => $operationData->expertId,
            'EXPERT_NAME' => $expert->fullName,
            'CLIENT_ID' => $client->id,
            'CLIENT_NAME' => $client->fullName,
            'CONSUMPTION_ID' => $operationData->consumptionId,
            'CONSUMPTION_NUMBER' => $operationData->consumptionNumber,
            'AGREEMENT_NUMBER' => $operationData->agreementNumber,
            'DATE' => $operationData->date,
            'DIFFERENCES' => $differences,
        ];
    }


}