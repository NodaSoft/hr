<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Contracts\ContractorServiceContract;
use NW\WebService\References\Operations\Notification\Contracts\EmployeeServiceContract;
use NW\WebService\References\Operations\Notification\Contracts\MessagesClientContract;
use NW\WebService\References\Operations\Notification\Contracts\NotificationManagerContract;
use NW\WebService\References\Operations\Notification\Contracts\ReturnOperationContract;
use NW\WebService\References\Operations\Notification\Contracts\SellerServiceContract;
use NW\WebService\References\Operations\Notification\Dto\EmailMessageDto;
use NW\WebService\References\Operations\Notification\Dto\ReturnOperationRequestDto;
use NW\WebService\References\Operations\Notification\Enums\NotificationEventEnum;
use NW\WebService\References\Operations\Notification\Enums\StatusEnum;
use NW\WebService\References\Operations\Notification\Exceptions\ReturnOperationException;
use Symfony\Component\HttpFoundation\Response;

/**
 * Класс ReturnOperation похоже создан для обработки возвратов, но кода который бы что-то делал
 * с самим возвратом не видно, а только рассылает уведомления в зависимости от статуса и входных данных
 *
 * Код, если убрать ошибки рабочий, но требует рефакторинга
 *
 * Много логики в одном методе класса (ReturnOperation::doOperation), разбил код на части
 *
 * У нас идет работа с реквестом, значит мы можем взять из него данные и уровнем выше (например в реквесте
 * laravel) провалидировать и наполнить какую нибудь DTO нужными данными
 * + есть странные места: $creator = Employee::getById((int)$data['creatorId']); , которые по хорошему
 * вообще не должны быть выполнены если входные данные для $data['creatorId'] <= 0
 *
 * Метод класса MessagesClient sendMessage принимает разные параметры, сложно понять как правильно,
 * т.к. нет описания класса, да и само использование статического метода не очень по этому заменил
 * на контракт MessagesClientContract
 *
 * Есть какой то странный $error, но видимо его объявление удалили
 *
 * Вместо использования классов типа Contractor напрямую следует использовать их контракты
 *
 * Метод Contractor::getFullName всегда будет возвращать строку, даже если name и id у Contractor будут null
 * метод вернет " ", а empty(" ") вернет false, получется все проверки пройдут
 *
 * Не информативные экспшены, если пробросится в лог весь реквест то можно будет еще что-то понять, если нет
 * то не понятно на каких данных произошла ошибка
 *
 * Я бы добавил еще логгирование
 */
class ReturnOperation extends ReferencesOperation implements ReturnOperationContract
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @param ReturnOperationRequestDto $requestDto
     * @param MessagesClientContract $messagesClient
     * @param NotificationManagerContract $notificationManager
     * @param SellerServiceContract $sellerService
     * @param ContractorServiceContract $contractorService
     * @param EmployeeServiceContract $employeeService
     * @param string $emailFrom
     * @param string[] $emails
     */
    public function __construct(
        private readonly ReturnOperationRequestDto $requestDto,
        private readonly MessagesClientContract $messagesClient,
        private readonly NotificationManagerContract $notificationManager,
        private readonly SellerServiceContract $sellerService,
        private readonly ContractorServiceContract $contractorService,
        private readonly EmployeeServiceContract $employeeService,
        private readonly string $emailFrom,
        private readonly array $emails,
    ) {
    }

    /**
     * {@inheritdoc}
     */
    public function doOperation(): array
    {
        $requestDto = $this->requestDto;
        $error = '';
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        $reseller = $this->sellerService->getById($requestDto->resellerId);

        if ($reseller === null) {
            throw new ReturnOperationException('Seller not found!', Response::HTTP_NOT_FOUND);
        }

        $client = $this->contractorService->getById($requestDto->clientId);

        if ($client === null) {
            throw new ReturnOperationException('Client not found!', Response::HTTP_NOT_FOUND);
        }

        if ($client->type !== Contractor::TYPE_CUSTOMER || $client->seller->id !== $requestDto->resellerId) {
            throw new ReturnOperationException('Invalid client!', Response::HTTP_INTERNAL_SERVER_ERROR);
        }

        $creator = $this->employeeService->getById($requestDto->creatorId);

        if ($creator === null) {
            throw new ReturnOperationException('Creator not found!', Response::HTTP_NOT_FOUND);
        }

        $expert = $this->employeeService->getById($requestDto->expertId);

        if ($expert === null) {
            throw new ReturnOperationException('Expert not found!', Response::HTTP_NOT_FOUND);
        }

        $additionalTemplateData = [
            'CREATOR_NAME' => $creator->getFullName(),
            'EXPERT_NAME' => $expert->getFullName(),
            'CLIENT_NAME' => empty($client->getFullName()) ? $client->name : $client->getFullName(),
            'DIFFERENCES' => $this->calculateDifference(
                $requestDto->notificationType,
                $requestDto->resellerId,
                $requestDto,
            ),
        ];

        $this->validateTemplateData($additionalTemplateData);

        $templateData = array_merge($this->fillTemplateData($requestDto), $additionalTemplateData);

        $result['notificationEmployeeByEmail'] = $this->sendEmailMessages(
            requestDto: $requestDto,
            emailFrom: $this->emailFrom,
            emails: $this->emails,
            subject: __('complaintEmployeeEmailSubject', $templateData, $requestDto->resellerId),
            message: __('complaintEmployeeEmailBody', $templateData, $requestDto->resellerId),
        );

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($requestDto->notificationType === self::TYPE_CHANGE && !empty($requestDto?->differences['to'])) {
            $result['notificationClientByEmail'] = $this->sendEmailMessages(
                requestDto: $requestDto,
                emailFrom: $this->emailFrom,
                emails: [$client->email],
                subject: __('complaintClientEmailSubject', $templateData, $requestDto->resellerId),
                message: __('complaintClientEmailBody', $templateData, $requestDto->resellerId),
                differencesTo: $requestDto->differences['to'],
            );

            if (!empty($client->mobile)) {
                $res = $this
                    ->notificationManager
                    ->send(
                        resellerId: $requestDto->resellerId,
                        clientId: $requestDto->clientId,
                        notificationType: NotificationEventEnum::Change->value,
                        differencesTo: $requestDto->differences['to'],
                        templateData: $templateData,
                    );

                if ($res) {
                    $result['notificationClientBySms']['isSent'] = true;
                }
                if (!empty($error)) {
                    $result['notificationClientBySms']['message'] = $error;
                }
            }
        }

        return $result;
    }

    /**
     * @throws ReturnOperationException
     */
    private function validateTemplateData(array $data): void
    {
        foreach ($data as $key => $tempData) {
            if (empty($tempData)) {
                throw new ReturnOperationException(
                    "Template Data ({$key}) is empty!",
                    Response::HTTP_INTERNAL_SERVER_ERROR
                );
            }
        }
    }

    /**
     * @throws ReturnOperationException
     */
    private function calculateDifference(
        int $notificationType,
        int $resellerId,
        ReturnOperationRequestDto $requestDto,
    ): string {
        if ($notificationType === self::TYPE_NEW) {
            return __('NewPositionAdded', null, $resellerId);
        }

        if ($notificationType === self::TYPE_CHANGE && $requestDto->differences !== null) {
            return __('PositionStatusHasChanged', [
                'FROM' => StatusEnum::from($requestDto->differences['from'])->name,
                'TO' => StatusEnum::from($requestDto->differences['to'])->name,
            ], $resellerId);
        }

        throw new ReturnOperationException('Empty difference!', Response::HTTP_BAD_REQUEST);
    }

    private function sendEmailMessages(
        ReturnOperationRequestDto $requestDto,
        string $emailFrom,
        array $emails,
        string $subject,
        string $message,
        int $differencesTo = null,
    ): bool {
        if (empty($this->emailFrom) || empty($this->emails)) {
            return false;
        }

        $result = false;

        foreach ($emails as $email) {
            if (empty($email)) {
                continue;
            }

            $this
                ->messagesClient
                ->sendMessages(
                    messages: [
                        new EmailMessageDto(
                            emailFrom: $emailFrom,
                            emailTo: $email,
                            subject: $subject,
                            message: $message,
                        )
                    ],
                    resellerId: $requestDto->resellerId,
                    clientId: $requestDto->clientId,
                    notificationType: $requestDto->notificationType,
                    differencesTo: $differencesTo,
                );

            $result = true;
        }

        return $result;
    }

    private function fillTemplateData(ReturnOperationRequestDto $requestDto): array
    {
        return [
            'COMPLAINT_ID' => $requestDto->complaintId,
            'COMPLAINT_NUMBER' => $requestDto->complaintNumber,
            'CREATOR_ID' => $requestDto->creatorId,
            'EXPERT_ID' => $requestDto->expertId,
            'CLIENT_ID' => $requestDto->clientId,
            'CONSUMPTION_ID' => $requestDto->consumptionId,
            'CONSUMPTION_NUMBER' => $requestDto->consumptionNumber,
            'AGREEMENT_NUMBER' => $requestDto->agreementNumber,
            'DATE' => $requestDto->date,
        ];
    }
}
