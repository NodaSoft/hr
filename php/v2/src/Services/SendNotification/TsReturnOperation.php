<?php

namespace Nodasoft\Testapp\Services\SendNotification;


use Exception;
use Nodasoft\Testapp\Config\MockConfigData;
use Nodasoft\Testapp\DTO\GetNotificationDifferenceDTO;
use Nodasoft\Testapp\DTO\MessageDifferenceDto;
use Nodasoft\Testapp\DTO\SendNotificationDTO;
use Nodasoft\Testapp\Entities\Client\Client;
use Nodasoft\Testapp\Entities\Employee\Employee;
use Nodasoft\Testapp\Entities\Seller\Seller;
use Nodasoft\Testapp\Enums\NotificationType;
use Nodasoft\Testapp\Events\Base\EventDispatcher;
use Nodasoft\Testapp\Events\ChangeReturnStatusEvent;
use Nodasoft\Testapp\Notifications\Base\NotificationDispatcher;
use Nodasoft\Testapp\Notifications\MailNotificationInterfaceClient;
use Nodasoft\Testapp\Notifications\SmsNotificationInterfaceClient;
use Nodasoft\Testapp\Repositories\Client\ClientRepositoryInterface;
use Nodasoft\Testapp\Repositories\Employee\EmployeeRepositoryInterface;
use Nodasoft\Testapp\Repositories\Seller\SellerRepositoryInterface;
use Nodasoft\Testapp\Services\SendNotification\Base\GetDifferencesInterface;
use Nodasoft\Testapp\Services\SendNotification\Base\ReferencesOperation;

class TsReturnOperation implements ReferencesOperation
{
    private Seller $reseller;
    private Client $client;
    private Employee $creator;
    private Employee $expert;
    private string $differences;
    private array $templateData;
    private string $fromEmail;

    public function __construct(
        private readonly ClientRepositoryInterface   $clientRepository,
        private readonly EmployeeRepositoryInterface $employeeRepository,
        private readonly SellerRepositoryInterface   $sellerRepository,
        private readonly GetDifferencesInterface     $differentService,
    )
    {
    }

    /**
     * @throws Exception
     */
    public function doOperation(SendNotificationDTO $dto): array
    {
        // Initialize result structure
        $result = [
            'notificationEmployeeByEmail' => [
                'isSet' => false,
                'message' => ''
            ],
            'notificationClientByEmail' => [
                'isSet' => false,
                'message' => ''
            ],
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        $this->reseller = $this->sellerRepository->getById($dto->resellerId);
        $this->client = $this->clientRepository->getById($dto->clientId);

        $this->ensureGivenClientIsCorrect();

        $this->creator = $this->employeeRepository->getById($dto->creatorId);
        $this->expert = $this->employeeRepository->getById($dto->expertId);

        // Determine differences message
        $this->differences = $this->getDifferences($dto->notificationType, $dto->differences);

        // Prepare template data
        $this->templateData = $this->prepareTemplate($dto);

        // Получаем email отправитля из настроек
        $this->fromEmail = MockConfigData::getResellerEmailFrom($this->reseller->getId());

        // Шлём уведомление сотрудинкам
        if ($this->fromEmail) {
            try {
                $this->notifyEmployeesByEmail();
                EventDispatcher::dispatch(new ChangeReturnStatusEvent([
                    'resellerId' => $this->reseller->getId()
                ]));
                $result['notificationEmployeeByEmail']['isSent'] = true;
            } catch (Exception $e) {
                $result['notificationEmployeeByEmail']['message'] = $e->getMessage();
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if (!($dto->notificationType->isChange() && $dto->differences)) return $result;

        if ($this->fromEmail && $this->client->getEmail()) {
            try {
                $this->notifyClientByEmail();
                $result['notificationClientByEmail']['isSent'] = true;
                EventDispatcher::dispatch(new ChangeReturnStatusEvent([
                    'reseller_id' => $this->reseller->getId(),
                    'client_id' => $this->client->getId(),
                    'status_id' => $dto->differences->to
                ]));
            } catch (Exception $e) {
                $result['notificationClientByEmail']['message'] = $e->getMessage();
            }
        }

        if ($this->client->getMobile()) {
            try {
                $this->notifyClientBySms();
                $result['notificationClientBySms']['isSent'] = true;
                EventDispatcher::dispatch(new ChangeReturnStatusEvent([
                    'reseller_id' => $this->reseller->getId(),
                    'client_id' => $this->client->getId(),
                    'status_id' => $dto->differences->to
                ]));
            } catch (Exception $e) {
                $result['notificationClientBySms']['message'] = $e->getMessage();
            }
        }

        return $result;
    }

    private function notifyEmployeesByEmail(): void
    {
        // Получаем email сотрудников из настроек
        $emails = MockConfigData::getEmailsByPermit($this->reseller->getId(), 'tsGoodsReturn');

        foreach ($emails as $to) {
            NotificationDispatcher::dispatch(new MailNotificationInterfaceClient([
                'emailFrom' => $this->fromEmail,
                'emailTo' => $to,
                'subject' => __('complaintEmployeeEmailSubject', $this->templateData, $this->reseller->getId()),
                'message' => __('complaintEmployeeEmailBody', $this->templateData, $this->reseller->getId()),
            ]));
        }
    }

    private function notifyClientByEmail(): void
    {
        NotificationDispatcher::dispatch(new MailNotificationInterfaceClient([
            'emailFrom' => $this->fromEmail,
            'emailTo' => $this->client->getEmail(),
            'subject' => __('complaintClientEmailSubject', $this->templateData, $this->reseller->getId()),
            'message' => __('complaintClientEmailBody', $this->templateData, $this->reseller->getId()),
        ]));

    }

    private function notifyClientBySms(): void
    {
        NotificationDispatcher::dispatch(new SmsNotificationInterfaceClient([
            'client_id' => $this->client->getId(),
            'mobile' => $this->client->getMobile(),
            'reseller' => $this->reseller->getId(),
            'template' => $this->templateData,
        ]));
    }

    private function getDifferences(NotificationType $type, MessageDifferenceDto $differenceDto)
    {
        return $this->differentService->getDifference(new GetNotificationDifferenceDTO(
            $type,
            $this->reseller->getId(),
            $differenceDto,
        ));
    }

    /**
     * @throws Exception
     */
    private function prepareTemplate(SendNotificationDTO $dto): array
    {
        $template = [
            'COMPLAINT_ID' => $dto->complaintId,
            'COMPLAINT_NUMBER' => $dto->complaintNumber,
            'CREATOR_ID' => $dto->creatorId,
            'CREATOR_NAME' => $this->creator->getFullName(),
            'EXPERT_ID' => $dto->expertId,
            'EXPERT_NAME' => $this->expert->getFullName(),
            'CLIENT_ID' => $dto->clientId,
            'CLIENT_NAME' => $this->client->getFullName(),
            'CONSUMPTION_ID' => $dto->consumptionId,
            'CONSUMPTION_NUMBER' => $dto->consumptionNumber,
            'AGREEMENT_NUMBER' => $dto->agreementNumber,
            'DATE' => $dto->date,
            'DIFFERENCES' => $this->differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($template as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ($key) is empty!", 500);
            }
        }

        return $template;
    }

    /**
     * @throws Exception
     */
    private function ensureGivenClientIsCorrect(): void
    {
        if (!$this->client->isCustomer() || !$this->client->HasSeller($this->reseller->getId()))
            throw new Exception('invalid client data!', 400);
    }
}

