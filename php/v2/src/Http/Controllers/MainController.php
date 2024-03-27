<?php

namespace App\Http\Controllers;

use App\Config\Config;
use App\Enum\ContractorType;
use App\Enum\Notification;
use App\Enum\Status;
use App\Exceptions\EntityNotFoundException;
use App\Http\Request\BaseRequest;
use App\Http\Request\NewPositionRequest;
use App\Http\Request\StatusChangedRequest;
use App\Models\Contractor;
use App\Models\Employee;
use App\Models\Seller;
use App\Services\MailService;
use App\Services\SmsService;

class MainController
{
    protected function handleBaseRequest(BaseRequest $request, array $diffs): array
    {
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        $resellerId = $request->getDataField(BaseRequest::RESELLER_ID);

        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new EntityNotFoundException('Reseller', $resellerId);
        }

        $clientId = $request->getDataField(BaseRequest::CLIENT_ID);
        $client = Contractor::getById($clientId);
        if (
            !$client instanceof Contractor ||
            $client->type !== ContractorType::CUSTOMER->value ||
            $client->seller->id !== $resellerId
        ) {
            throw new EntityNotFoundException('Client', $clientId);
        }

        $employeeId = $request->getDataField(BaseRequest::CREATOR_ID);
        $creator = Employee::getById($employeeId);
        if (!$creator instanceof Employee) {
            throw new EntityNotFoundException('Creator', $employeeId);
        }

        $expertId = $request->getDataField(BaseRequest::EXPERT_ID);
        $expert = Employee::getById($expertId);
        if ($expert === null) {
            throw new EntityNotFoundException('Expert', $employeeId);
        }

        $templateData = $this->getTemplate($request, $creator, $expertId, $client, $diffs);

        $emailFrom = Config::RESELLER_FROM_EMAIL;
        // Получаем email сотрудников из настроек
        $emails = Config::getEmailsByEventForReseller($resellerId, 'tsGoodsReturn');
        if (count($emails) > 0) {
            $result['notificationEmployeeByEmail'] = true;
            foreach ($emails as $email) {
                MailService::sendMessage(
                    [
                    'emailFrom' => $emailFrom,
                    'emailTo' => $email,
                    'subject' => MailService::getSubject($templateData, $resellerId),
                    'message' => MailService::getMessage($templateData, $resellerId),
                ],
                    $resellerId,
                    Config::EVENT_CHANGE_RETURN_STATUS
                );
            }
        }

        if ($client->email !== null) {
            $result['notificationClientByEmail'] = true;

            MailService::sendMessage([
                'emailFrom' => $emailFrom,
                'emailTo' => $client->email,
                'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                'message' => __('complaintClientEmailBody', $templateData, $resellerId),
            ], $resellerId, Config::EVENT_CHANGE_RETURN_STATUS, [ $client->id, $diffs]);
        }

        if ($client->mobile !== null) {
            [$res, $error] = SmsService::sendMessage([
                $resellerId,
                $client->id,
                Config::EVENT_CHANGE_RETURN_STATUS,
                $diffs,
                $templateData]);
            if ($res) {
                $result['notificationClientBySms']['isSent'] = true;
            }
            if (!empty($error)) {
                $result['notificationClientBySms']['message'] = $error;
            }
        }

        return $result;
    }

    protected function calculateDiffs(string $type, int $resellerId, ?array $options = null): array
    {
        //do something
        return [];
    }

    protected function getTemplate(
        BaseRequest $request,
        Employee $creator,
        Employee $expert,
        Contractor $client,
        array $differences
    ) {
        return [
            'COMPLAINT_ID' => $request->getDataField(BaseRequest::COMPLAINT)[BaseRequest::ID],
            'COMPLAINT_NUMBER' => $request->getDataField(BaseRequest::COMPLAINT)[BaseRequest::NUMBER],
            'CREATOR_ID' => $creator->id,
            'CREATOR_NAME' => $creator->getFullName(),
            'EXPERT_ID' => $expert->id,
            'EXPERT_NAME' => $expert->getFullName(),
            'CLIENT_ID' => $client->id,
            'CLIENT_NAME' => $client->getFullName(),
            'CONSUMPTION_ID' => $request->getDataField(BaseRequest::CONSUMPTION)[BaseRequest::ID],
            'CONSUMPTION_NUMBER' => $request->getDataField(BaseRequest::CONSUMPTION)[BaseRequest::NUMBER],
            'AGREEMENT_NUMBER' => $request->getDataField(BaseRequest::AGREEMENT_NUMBER),
            'DATE' => $request->getDataField(BaseRequest::DATE),
            'DIFFERENCES' => $differences,
        ];
    }

    /**
     * @throws EntityNotFoundException
     */
    public function newPosition(NewPositionRequest $request): array
    {
        $resellerId = $request->getDataField(BaseRequest::RESELLER_ID);

        $differences = $this->calculateDiffs('NewPositionAdded', $resellerId);
        return $this->handleBaseRequest($request, $differences);
    }

    public function statusChanged(StatusChangedRequest $request)
    {
        $resellerId = $request->getDataField(BaseRequest::RESELLER_ID);

        $differences = $this->calculateDiffs('PositionStatusHasChanged', $resellerId, [
            'TO' => Status::getNameById($request->getDataField(StatusChangedRequest::NEW_STATUS)),
        ]);

        return $this->handleBaseRequest($request, $differences);
    }
}
