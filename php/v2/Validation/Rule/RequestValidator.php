<?php

namespace NW\WebService\References\Operations\Notification\Validation\Rule;

use NW\WebService\References\Operations\Notification\Validation\ValidatorInterface;

/**
 * RequestValidator class
 */
class RequestValidator implements ValidatorInterface
{

    public function validate(array $data, array &$result = []): bool {
        if (empty((int)($data['resellerId']))) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return false;
        }

        $requiredFieldsErrors = [
            'clientId' => 'Empty clientId',
            'creatorId' => 'Empty creatorId',
            'expertId' => 'Empty expertId',
            'complaintId' => 'Empty complaintId',
            'complaintNumber' => 'Empty complaintNumber',
            'consumptionId' => 'Empty consumptionId',
            'consumptionNumber' => 'Empty consumptionNumber',
            'agreementNumber' => 'Empty agreementNumber',
            'date' => 'Empty date',
            'notificationType' => 'Empty notificationType'
        ];

        foreach ($requiredFieldsErrors as $field => $errorMessage) {
            if (empty(trim((string)($data[$field] ?? '')))) {
                throw new \Exception($errorMessage, 400);
            }
        }

        return true;
    }
}
