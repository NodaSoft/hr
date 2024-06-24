<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\ReferencesOperation;
use NW\Models\Contractor;
use NW\Models\Employee;
use NW\Models\Seller;
use NW\Models\Employee;
use NW\Models\Contractor;
use NW\Services\Status;
use NW\Services\MessagesClient;
use stdClass;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;


    private array $data, $result;
    private $constructData;

    public function __construct()
    {
        $this->constructData = new stdClass();
        $this->data = (array)$this->getRequest("data");
        $this->result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];
    }


    public function doOperation()
    {
        // Упорядочил по нагрузке, чтобы не было лишних запросов к бд
        try {
            $this->handlerComplaint(); // +

            $this->handlerNotification(); //+

            $this->handlerResseller(); // +

            $this->handlerDifferences(); // +

            $this->handlerClient(); // +

            $this->handlerEmployee(); // +

            $this->handlerConsumption(); // +

            $this->handlerAgreement(); // +

            $templateData = $this->handlerTemplate(); // +

            $this->sendMessageEmail($templateData); // +

            $this->handlerResult($templateData); // +

            return $this->result;
        } catch (\Exception $e) {
            $errorMessage = $e->getMessage() ?? "Unknown error";
            throw new \Exception($errorMessage);
        }
    }

    public function handlerAgreement()
    {
        $agreementNumber = isset($this->data) ? (string)$this->data['agreementNumber'] : null;

        if (empty($agreementNumber))
            throw new \Exception("Agreement number not found", 400);

        $agreement = new stdClass();
        $agreement->number = $agreementNumber;

        $this->constructData->agreement = $agreement;
    }

    public function handlerConsumption()
    {
        $consumptionId = isset($this->data) ? (int)$this->data['consumptionId'] : 0;
        $consumptionNumber = isset($this->data) ? (string)$this->data['consumptionNumber'] : 0;

        if (empty($consumptionId))
            throw new \Exception("Consumption id not found", 400);

        $consumption = new stdClass();
        $consumption->id = $consumptionId;
        $consumption->number = $consumptionNumber;

        $this->constructData->consumption = $consumption;
    }


    public function handlerResult($templateData)
    {
        $notificationType = $this->constructData->notificationType;
        $client = $this->constructData->client;
        $resellerId = $this->constructData->reseller->id;
        $differencesTo = isset($this->data['differences']['to']) ? (int) $this->data['differences']['to'] : 0;
        if ($notificationType === self::TYPE_CHANGE && !empty($differencesTo)) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $client->email,
                        'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, ($differencesTo));
                $this->result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                try {
                    $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo, $templateData);
                    if ($res)
                        $this->result['notificationClientBySms']['isSent'] = true;
                } catch (\Exception $e) {
                    $errorMessage = $e->getMessage() ?? "Unknown error";
                    $this->result['notificationClientBySms']['message'] = $errorMessage;
                }
            }
        }
    }

    public function sendMessageEmail($templateData)
    {
        $resellerId = $this->constructData->reseller->id ?? null;

        $emailFrom = Seller::getResellerEmailFrom($resellerId);
        $emails = Seller::getEmailsByPermit($resellerId, 'tsGoodsReturn');

        if (empty($emailFrom))
            throw new \Exception("Email from is empty");

        if (empty($emails))
            throw new \Exception("Emails is empty");

        foreach ($emails as $email) {
            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $email,
                    'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
            $this->result['notificationEmployeeByEmail'] = true;
        }
    }

    // Поскольку передавать более 3х аргументов не принято, поэтому сделал свойством, возможно здесь и было бы уместно более 3х поставить, но не будем.
    public function handlerTemplate()
    {
        $complaintId = (int)$this->constructData->complaint->id ?? null;
        $complaintNumber = (string)$this->constructData->complaint->number ?? null;
        $creatorId = (int)$this->constructData->employees->creator->id ?? null;
        $expertId = (int)$this->constructData->employees->expert->id ?? null;
        $clientId = (int)$this->constructData->client->id ?? null;
        $consumptionId = (int) $this->constructData->consumption->id ?? null;
        $consumptionNumber = (string) $this->constructData->consumption->number ?? null;
        $agreementNumber = (string)$this->constructData->agreement->number ?? null;
        $differences = (string)$this->constructData->differences;
        $date = (string)$this->constructData->date;

        // Я не могу с уверенностью сказать как у вас сделано, чтобы оставить объект ФИО, не думаю что вы перебираете массив или сериализуете, поэтому вставляю только имя, да и ключ в темплейте NAME.
        $creatorName = (string) $this->constructData->employees->creator->name ?? null;
        $expertName = (string) $this->constructData->employees->expert->name ?? null;
        $clientName = (string) $this->constructData->client->name ?? null;

        $templateData = [
            'COMPLAINT_ID'       => $complaintId,
            'COMPLAINT_NUMBER'   => $complaintNumber,
            'CREATOR_ID'         => $creatorId,
            'CREATOR_NAME'       => $creatorName,
            'EXPERT_ID'          => $expertId,
            'EXPERT_NAME'        => $expertName,
            'CLIENT_ID'          => $clientId,
            'CLIENT_NAME'        => $clientName,
            'CONSUMPTION_ID'     => $consumptionId,
            'CONSUMPTION_NUMBER' => $consumptionNumber,
            'AGREEMENT_NUMBER'   => $agreementNumber,
            'DATE'               => $date,
            'DIFFERENCES'        => $differences,
        ];


        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData))
                throw new \Exception("Template Data ({$key}) is empty!", 500);
        }
        return $templateData;
    }


    public function handlerComplaint(): object
    {
        $complaintId = isset($this->data['complaintId']) ? (int) $this->data['complaintId'] : null;
        // возможно здесь необходим тип данных Integer, т.к присутствует Number в названии переменной, но думаю у вас в формате 1gdf325-5235j325hg-5115, поэтому оставлю так
        $complaintNumber = isset($this->data['complaintNumber']) ? (string) $this->data['complaintNumber'] : null;

        if (empty($complaintId) || empty($complaintNumber))
            throw new \Exception("Complaint not full field", 400);

        $complaint = new stdClass();
        $complaint->id = $complaintId;
        $complaint->number = $complaintNumber;
        return $complaint;
    }

    public function handlerDifferences()
    {
        $dataDifferences = isset($data['differences']) ? $data['differences'] : null;
        $resellerId = $this->constructData->reseller->id ?? 0;
        $notificationType = $this->constructData->notificationType;

        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($dataDifferences)) {
            $differencesFrom = isset($data['differences']['from']) ? (int)$data['differences']['from'] : 0;
            $differencesTo = isset($data['differences']['to']) ? (int)$data['differences']['to'] : 0;

            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName($differencesFrom),
                'TO'   => Status::getName($differencesTo),
            ], $resellerId);
        }

        if (!isset($differences))
            throw new \Exception("Differences is empty");

        $this->constructData->differences = $differences;
    }

    public function handlerEmployee()
    {

        // Проверку на ФИО быть может делать было необязательно

        // Быть может надо было разъединить, но можно и так.
        $creatorId = isset($data['creatorId']) ? (int) $data['creatorId'] : null;
        $expertId = isset($data['expertId']) ? (int) $data['expertId'] : null;

        if (empty($creatorId))
            throw new \Exception("Creator not found", 400);

        if (empty($expertId))
            throw new \Exception("Expert not found", 400);


        $creator = Employee::getById($creatorId);
        if (empty($creator))
            throw new \Exception('Creator not found!', 400);

        $creatorFullName = $creator->getFullName();
        $creatorName = !empty($creatorFullName) && isset($creatorFullName->name) ? $creatorFullName->name : null;
        if (empty($creatorName))
            throw new \Exception("Creator not has name", 400);

        $creator->name = $creatorName;
        // --
        $expert = Employee::getById($expertId);
        if ($expert === null)
            throw new \Exception('Expert not found!', 400);

        $expertFullName = $expert->getFullName();
        $expertName = !empty($expertFullName) && isset($expertFullName->name) ? $expertFullName->name : null;
        if (empty($expertName))
            throw new \Exception("Expert not has name", 400);

        $expert->name = $expertName;

        $employees = new stdClass();
        $employees->creator = $creator;
        $employees->expert = $expert;
        return $employees;
    }

    public function handlerClient(): void
    {
        $clientId = isset($this->data['clientId']) ? (int) $this->data['clientId'] : null;
        $resellerId = $this->constructData->reseller->id;

        if (empty($clientId))
            throw new \Exception('Client not found!', 400);

        $client = Contractor::getById($clientId);

        // возможно здесь надо $client->id !== $resellerId. но я не знаю как у вас устроено
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Client not found!', 400);
        }

        $clientFullName = $client->getFullName();
        $clientName = !empty($clientFullName) && isset($client->name) ? $client->name : null;

        if (empty($clientName))
            throw new \Exception("Client not has name", 400);

        $client->name = $clientName;

        $this->constructData->client = $client;
    }

    public function handlerResseller(): void
    {
        $resellerId = isset($this->data['resellerId']) ? (int)$this->data['resellerId'] : null;

        if (empty($resellerId))
            throw new \Exception("Seller not found", 400);
        // $this->result['notificationClientBySms']['message'] = 'Empty resellerId';


        $reseller = Seller::getById($resellerId);
        if ($reseller === null)
            throw new \Exception('Seller not found!', 400);

        $this->constructData->reseller = $reseller;
    }

    public function handlerNotification(): void
    {
        $notificationType = isset($data['notificationType']) ? (int) $data['notificationType'] : null;

        if (empty($notificationType))
            throw new \Exception('Empty notificationType', 400);
        $this->constructData->notificationType = $notificationType;
    }
}
