<?php
/**
 * @author    Vasyl Minikh <mhbasil1@gmail.com>
 * @copyright 2024
 *
 */
declare(strict_types=1);

use NW\WebService\References\Operations\Notification\Exception\ClientNotFoundException;
use NW\WebService\References\Operations\Notification\Exception\CreatorNotFoundException;
use NW\WebService\References\Operations\Notification\Exception\EmptyNotificationTypeException;
use NW\WebService\References\Operations\Notification\Exception\ExpertNotFoundException;
use NW\WebService\References\Operations\Notification\Exception\SelerNotFoundException;
use NW\WebService\References\Operations\Notification\Contracts\RequestOperationAbstract;
use NW\WebService\References\Operations\Notification\Contracts\Contractor;
use NW\WebService\References\Operations\Notification\Contracts\Seller;
use NW\WebService\References\Operations\Notification\Contracts\Employee;
use NW\WebService\References\Operations\Notification\Contracts\Status;
use NW\WebService\References\Operations\Notification\Dto\OperationResponseDto;
use NW\WebService\References\Operations\Notification\Dto\OperationRequestDto;
use NW\WebService\References\Operations\Notification\Dto\EmailMessageDto;
use NW\WebService\References\Operations\Notification\Contracts\MessagesClient;
use NW\WebService\References\Operations\Notification\Contracts\NotificationManager;

/**
 * Class ReturnOperation.
 *
 */
class ReturnOperation extends RequestOperationAbstract
{
    private const TYPE_NEW = 1;
    private const TYPE_CHANGE = 2;
    /**
     * @var Contractor
     */
    private Contractor $contractor;
    /**
     * @var Seller
     */
    private Seller $seller;
    /**
     * @var Employee
     */
    private Employee $employee;
    /**
     * @var Status
     */
    private Status $status;
    /**
     * @var OperationResponseDto
     */
    private OperationResponseDto $response;
    /**
     * @var OperationRequestDto
     */
    private OperationRequestDto $request;
    /**
     * @var MessagesClient
     */
    private MessagesClient $messagesClient;
    /**
     * @var NotificationManager
     */
    private NotificationManager $notificationManager;

    /**
     * ReturnOperation constructor.
     *
     * @param Contractor $contractor
     * @param Seller $seller
     * @param Employee $employee
     * @param Status $status
     * @param OperationResponseDto $response
     * @param OperationRequestDto $request
     * @param MessagesClient $messagesClient
     * @param NotificationManager $notificationManager
     */
    public function __construct(
        Contractor           $contractor,
        Seller               $seller,
        Employee             $employee,
        Status               $status,
        OperationResponseDto $response,
        OperationRequestDto  $request,
        MessagesClient       $messagesClient,
        NotificationManager  $notificationManager
    )
    {
        $this->contractor = $contractor;
        $this->seller = $seller;
        $this->employee = $employee;
        $this->status = $status;
        $this->response = $response;
        $this->request = $request;
        $this->messagesClient = $messagesClient;
        $this->notificationManager = $notificationManager;
    }

    /**
     * validate Request
     *
     * @return array|null
     *
     * @throws Exception
     */
    private function validateRequest(): ?array
    {
        $result = $this->response;
        if (!($this->getRequest('data') ?? false)) {
            return $result->toArray();
        }
        if ($this->request->getResellerId() == null) {
            $result->setMessage('Empty resellerId');
            return $result->toArray();
        }
        if ($this->request->getNotificationType() === null) {
            throw new EmptyNotificationTypeException();
        }
        $resellerId = $this->request->getResellerId();
        if ($this->seller->getById($resellerId) === null) {
            // смысла в проверке нет, т.к. ресселер всегда будет заполнен,
            // но нужно оставить, т.к. инициализация может быть другая
            throw new SelerNotFoundException();
        }

        if ($this->contractor->getById($this->request->getClientId()) === null ||
            $this->contractor->getType() !== Contractor::TYPE_CUSTOMER || $this->contractor->getId() !== $resellerId) {
            // смысла в проверке нет, т.к. клиент всегда будет заполнен,
            // но нужно оставить, т.к. инициализация может быть другая
            throw new ClientNotFoundException();
        }

        if ($this->employee->getById($this->request->getCreatorId()) === null) {
            // смысла в проверке нет, т.к. криэйтор всегда будет заполнен,
            // но нужно оставить, т.к. инициализация может быть другая
            throw new CreatorNotFoundException();
        }

        if ($this->employee->getById($this->request->getExpertId()) === null) {
            // смысла в проверке нет, т.к. криэйтор всегда будет заполнен,
            // но нужно оставить, т.к. инициализация может быть другая
            throw new ExpertNotFoundException();
        }

        return null;
    }

    /**
     * get Differences
     *
     * @return string
     *
     */
    private function getDifferences(): string
    {
        $notificationType = $this->request->getNotificationType();
        $resellerId = $this->request->getResellerId();
        $differences = $this->request->getDifferences();
        if ($notificationType === self::TYPE_NEW) {
            return __('NewPositionAdded', null, $resellerId);
        }
        if ($notificationType === self::TYPE_CHANGE && !empty($this->request->getDifferences())) {
            return __('PositionStatusHasChanged', [
                'FROM' => $this->status->getName((int)$differences['from']),
                'TO' => $this->status->getName((int)$differences['to']),
            ], $resellerId);
        }
        return '';
    }

    /**
     * Execute
     *
     * @return array
     *
     * @throws Exception
     */
    public function doOperation(): array
    {
        $result = $this->response;
        $data = (array)$this->getRequest('data');

        $this->request->filDto($data);

        $validateRequest = $this->validateRequest();

        if ($validateRequest !== null) {
            return $validateRequest;
        }


        $resellerId = $this->request->getResellerId();
        $reseller = $this->seller->getById($resellerId);
        $client = $this->contractor->getById($this->request->getClientId());
        $cr = $this->employee->getById($this->request->getCreatorId());
        $et = $this->employee->getById($this->request->getExpertId());
        $differences = $this->getDifferences();

        $templateData = [
            'COMPLAINT_ID' => $this->request->getComplaintId(),
            'COMPLAINT_NUMBER' => $this->request->getComplaintNumber(),
            'CREATOR_ID' => $this->request->getCreatorId(),
            'CREATOR_NAME' => $cr->getFullName(),
            'EXPERT_ID' => $this->request->getExpertId(),
            'EXPERT_NAME' => $et->getFullName(),
            'CLIENT_ID' => $this->request->getClientId(),
            'CLIENT_NAME' => $client->getFullName(),
            'CONSUMPTION_ID' => $this->request->getConsumptionId(),
            'CONSUMPTION_NUMBER' => $this->request->getConsumptionNumber(),
            'AGREEMENT_NUMBER' => $this->request->getAgreementNumber(),
            'DATE' => $this->request->getDate(),
            'DIFFERENCES' => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = $this->messagesClient->getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = $this->messagesClient->getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                $emailMess = new EmailMessageDto(
                    $emailFrom,
                    $email,
                    __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    __('complaintEmployeeEmailBody', $templateData, $resellerId)
                );
                $this->messagesClient->sendMessages(
                    [$emailMess->toArray()],
                    $resellerId,
                    null, // при null шлем только ресселеру
                    NotificationManager::CHANGE_RETURN_STATUS
                );
                $result->setNotificationEmployeeByEmail();

            }
        }
        $requestDifferences = $this->request->getDifferences();
        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($this->request->getNotificationType() === self::TYPE_CHANGE && !empty($requestDifferences['to'])) {
            if (!empty($emailFrom) && !empty($client->getEmail())) {
                $emailMess = new EmailMessageDto(
                    $emailFrom,
                    $client->getEmail(),
                    __('complaintClientEmailSubject', $templateData, $resellerId),
                    __('complaintClientEmailBody', $templateData, $resellerId)
                );
                $this->messagesClient->sendMessages(
                    [$emailMess->toArray()],
                    $resellerId,
                    $client->getId(),
                    NotificationManager::CHANGE_RETURN_STATUS,
                    (int)$requestDifferences['to']
                );
                $result['notificationClientByEmail'] = true;
            }

            $error = null; //потерялась инициализация переменой

            if ($client->isMobile()) {
                if (
                    $this->notificationManager->send(
                        $resellerId,
                        $client->getId(),
                        NotificationManager::CHANGE_RETURN_STATUS,
                        (int)$requestDifferences['to'],
                        $templateData,
                        $error
                    )
                ) {
                    $result->setIsSent();
                }
                if (!empty($error)) {
                    $result->setMessage($error);
                }
            }
        }

        return $result->toArray();
    }
}
