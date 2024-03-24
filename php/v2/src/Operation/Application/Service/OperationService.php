<?php

namespace Src\Operation\Application\Service;

use Src\Notification\Infrastructure\API\NotificationApi;
use Src\Operation\Application\Action\GetTemplateDataAction;
use Src\Operation\Application\Action\ValidateClientAction;
use Src\Operation\Application\DataTransferObject\OperationData;
use src\Operation\Application\Exceptions\ClientNotFoundException;
use src\Operation\Application\Exceptions\ContractorNotFoundException;
use src\Operation\Application\Exceptions\EmployeeNotFoundException;
use src\Operation\Application\Exceptions\SellerNotFoundException;
use src\Operation\Domain\Enum\NotificationType;
use src\Operation\Infrastructure\Adapters\ContractorAdapter;
use src\Operation\Infrastructure\Adapters\EmployeeAdapter;
use src\Operation\Infrastructure\Adapters\SellerAdapter;

final class OperationService
{
    private ContractorAdapter $contractorAdapter;
    private EmployeeAdapter $employeeAdapter;

    private SellerAdapter $sellerAdapter;

    private NotificationApi $notificationAdapter;

    public function __construct()
    {
        $this->contractorAdapter = new ContractorAdapter();
        $this->employeeAdapter = new EmployeeAdapter();
        $this->sellerAdapter = new SellerAdapter();
        $this->notificationAdapter = new NotificationApi();
    }

    /**
     * @param OperationData $operationData
     * @return array
     * @throws ClientNotFoundException
     * @throws ContractorNotFoundException
     * @throws EmployeeNotFoundException
     * @throws SellerNotFoundException
     */
    public function sendReturnNotification(OperationData $operationData): array
    {
        try {
            $client = $this->contractorAdapter->getById($operationData->clientId);
            $emails = $this->contractorAdapter->getEmailsByPermit($operationData->resellerId, 'operation_return');
            $reseller = $this->sellerAdapter->getById($operationData->resellerId);
            $creator = $this->employeeAdapter->getById($operationData->creatorId);
            $expert = $this->employeeAdapter->getById($operationData->expertId);

            (new ValidateClientAction)->execute($client, $reseller);

            $templateData = (new GetTemplateDataAction)->execute(
                $operationData,
                $reseller,
                $creator,
                $expert,
                $client
            );

            $this->notificationAdapter->sendEmailNotification([
                'emailFrom' => env('EMAIL_FROM'),
                'emails' => $emails,
                'templateData' => $templateData,
                'resellerId' => $reseller->id
            ]);

            // Шлём клиентское уведомление, только если произошла смена статуса
            if ($operationData->notificationType === NotificationType::TYPE_CHANGE
                && !empty($operationData->differences->to)
                && !empty($client->email)
            ) {

                $this->notificationAdapter->sendEmailNotification([
                    'emailFrom' => env('EMAIL_FROM'),
                    'emails' => [$client->email],
                    'templateData' => $templateData,
                    'resellerId' => $reseller->id
                ]);

            }

            if (!empty($client->mobile)) {
                $this->notificationAdapter->sendSmsNotification(
                    [
                        'phoneNumber' => $client,
                        'message' => $templateData,
                        'resellerId' => $reseller->id,
                        'clientId' => $operationData->clientId,
                        'status' => $operationData->notificationType
                    ]
                );
            }

            return ['success' => 'Notification sent'];

        } catch (ClientNotFoundException|ContractorNotFoundException|EmployeeNotFoundException|SellerNotFoundException $e) {
            // Логируем ошибку в зависимости от типа исключения
            // Отправляем уведомление в контроллер
            return ['error' => $e->getMessage()];
        }
    }

}
