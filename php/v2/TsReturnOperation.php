<?php

namespace NW\WebService\References\Operations\Notification;


use Contractor;
use Employee;
use Exception;
use NotificationEvents;
use NW\WebService\References\Operations\Notification\Contracts\MessagesClientInterface;
use NW\WebService\References\Operations\Notification\Contracts\NotificationManagerInterface;
use ReferencesOperation;
use Seller;
use Status;

/*
 * Проблемы исходного кода
 *
 * Ошибка синтаксиса в возвращаемом типе метода doOperation:
 *    В начале метода doOperation был указан возвращаемый тип void, однако фактически метод возвращал массив с данными.
 *    Это приводило к ошибке синтаксиса.
 *
 * Жесткая связь с внешними сервисами:
 *    Класс TsReturnOperation напрямую вызывал методы для отправки сообщений и уведомлений
 *    (MessagesClient::sendMessage и NotificationManager::send). Это затрудняло тестирование,
 *    так как для тестов необходимо было бы эмулировать поведение этих сервисов.
 *
 * Отсутствие DTO (Data Transfer Object):
 *    Данные из запроса извлекались и использовались напрямую. Это увеличивало вероятность ошибок
 *    и затрудняло понимание структуры данных, передаваемых в метод.
 *
 * Отсутствие валидации и проверок в отдельных методах:
 *    Весь процесс валидации и обработки данных был сосредоточен в одном большом методе doOperation,
 *    что делало его сложным для понимания и поддержки.
 *
 * Смешение логики представления и обработки данных:
 *    Метод doOperation занимался как обработкой данных, так и формированием сообщений и их отправкой.
 *    Это нарушало принцип единственной ответственности.
 *
 * Использование явной проверки на null:
 *    В некоторых местах кода использовалась явная проверка на null с помощью условия вида `$client === null`.
 *    Это усложняет читаемость кода и может приводить к ошибкам из-за неправильного использования.
 *
 * Внесенные изменения
 *
 * Интерфейсы для внешних сервисов:
 *    Были созданы интерфейсы MessagesClientInterface и NotificationManagerInterface, что позволило
 *    абстрагировать реализацию отправки сообщений и уведомлений от самого класса TsReturnOperation.
 *
 * Внедрение зависимостей через конструктор:
 *    Реализации интерфейсов были переданы в класс через конструктор. Это облегчило тестирование с помощью мок-объектов.
 *
 * Использование DTO:
 *    Были созданы классы TsReturnOperationRequest и TsReturnOperationResponse для четкого определения
 *    структуры входных и выходных данных. Это улучшило читаемость и поддержку кода.
 *
 * Разделение логики на отдельные методы:
 *    Валидация данных и обработка были разделены на несколько методов. Это упростило основной метод
 *    doOperation и улучшило читаемость кода.
 *
 *
 * Итоговый результат
 *
 * Улучшенная тестируемость:
 *    Использование интерфейсов и инъекция зависимостей позволили легко создавать моки для тестирования,
 *    что упростило написание и выполнение unit-тестов.
 *
 * Улучшенная поддерживаемость:
 *    Разделение логики на несколько методов и использование DTO сделали код более организованным и легким для поддержки.
 *
 * Улучшенная читаемость и устойчивость к ошибкам:
 *    Явное определение структуры данных и отдельных методов для валидации и обработки данных улучшили читаемость кода
 *    и уменьшили вероятность ошибок.
 *
 * Заключение
 * В ходе рефакторинга я улучшил архитектуру исходного кода, применив принципы SOLID и улучшив его модульность.
 * Это позволило сделать код более гибким, легко тестируемым и поддерживаемым. Теперь класс TsReturnOperation не зависит
 * напрямую от реализации внешних сервисов и четко определяет структуру входных и выходных данных, что делает его более
 * устойчивым к изменениям и ошибкам.
 */


class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    private MessagesClientInterface $messagesClient;
    private NotificationManagerInterface $notificationManager;

    public function __construct(MessagesClientInterface $messagesClient, NotificationManagerInterface $notificationManager)
    {
        $this->messagesClient = $messagesClient;
        $this->notificationManager = $notificationManager;
    }

    /**
     * @throws Exception
     */
    public function doOperation(): array
    {
        $data = new TsReturnOperationRequest($this->getRequest('data'));
        $response = new TsReturnOperationResponse();

        $this->validateRequest($data);

        $reseller = Seller::getById($data->resellerId);
        $client = $this->validateClient($data->clientId, $reseller->id);
        $creator = $this->validateEmployee($data->creatorId, 'Creator');
        $expert = $this->validateEmployee($data->expertId, 'Expert');

        $templateData = $this->prepareTemplateData($data, $client, $creator, $expert);

        $this->sendEmployeeNotifications($templateData, $reseller->id, $response);
        $this->sendClientNotifications($data, $templateData, $client, $reseller->id, $response);

        return $response->toArray();
    }

    /**
     * @throws Exception
     */
    private function validateRequest(TsReturnOperationRequest $data): void
    {
        if (empty($data->resellerId)) {
            throw new Exception('Empty resellerId', 400);
        }

        if (empty($data->notificationType)) {
            throw new Exception('Empty notificationType', 400);
        }
    }

    private function validateClient(int $clientId, int $resellerId): Contractor
    {
        $client = Contractor::getById($clientId);
        if (!$client || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new Exception('Client not found!', 400);
        }
        return $client;
    }

    private function validateEmployee(int $employeeId, string $role): Employee
    {
        $employee = Employee::getById($employeeId);
        if (!$employee) {
            throw new Exception("{$role} not found!", 400);
        }
        return $employee;
    }

    private function prepareTemplateData(TsReturnOperationRequest $data, Contractor $client, Employee $creator, Employee $expert): array
    {
        $cFullName = $client->getFullName();
        if (empty($cFullName)) {
            $cFullName = $client->name;
        }

        $differences = '';
        if ($data->notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $data->resellerId);
        } elseif ($data->notificationType === self::TYPE_CHANGE && !empty($data->differences)) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data->differences['from']),
                'TO' => Status::getName((int)$data->differences['to']),
            ], $data->resellerId);
        }

        return [
            'COMPLAINT_ID' => $data->complaintId,
            'COMPLAINT_NUMBER' => $data->complaintNumber,
            'CREATOR_ID' => $data->creatorId,
            'CREATOR_NAME' => $creator->getFullName(),
            'EXPERT_ID' => $data->expertId,
            'EXPERT_NAME' => $expert->getFullName(),
            'CLIENT_ID' => $data->clientId,
            'CLIENT_NAME' => $cFullName,
            'CONSUMPTION_ID' => $data->consumptionId,
            'CONSUMPTION_NUMBER' => $data->consumptionNumber,
            'AGREEMENT_NUMBER' => $data->agreementNumber,
            'DATE' => $data->date,
            'DIFFERENCES' => $differences,
        ];
    }

    private function sendEmployeeNotifications(array $templateData, int $resellerId, TsReturnOperationResponse $response): void
    {
        $emailFrom = getResellerEmailFrom($resellerId);
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                $this->messagesClient->sendMessage([
                    0 => [
                        'emailFrom' => $emailFrom,
                        'emailTo' => $email,
                        'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $response->notificationEmployeeByEmail = true;
            }
        }
    }

    private function sendClientNotifications(TsReturnOperationRequest $data, array $templateData, Contractor $client, int $resellerId, TsReturnOperationResponse $response): void
    {
        if ($data->notificationType === self::TYPE_CHANGE && !empty($data->differences['to'])) {
            $emailFrom = getResellerEmailFrom($resellerId);

            if (!empty($emailFrom) && !empty($client->email)) {
                $this->messagesClient->sendMessage([
                    0 => [
                        'emailFrom' => $emailFrom,
                        'emailTo' => $client->email,
                        'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data->differences['to']);
                $response->notificationClientByEmail = true;
            }

            if (!empty($client->mobile)) {
                $res = $this->notificationManager->send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data->differences['to'], $templateData, $error);
                if ($res) {
                    $response->notificationClientBySms['isSent'] = true;
                }
                if (!empty($error)) {
                    $response->notificationClientBySms['message'] = $error;
                }
            }
        }
    }
}

