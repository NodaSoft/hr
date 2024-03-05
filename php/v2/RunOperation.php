<?php

namespace NW\WebService;


use NW\WebService\Config\Emails;
use NW\WebService\Employees\Employee;
use NW\WebService\Employees\EmployeeTypeEnum;
use NW\WebService\Exceptions\ValidationException;
use NW\WebService\Notification\NotificationTypeEnum;
use NW\WebService\Notifications\NotificationsEvents;
use NW\WebService\Request\RequestMapper;
use NW\WebService\Response\ResultOperation;
use NW\WebService\Validation\ValidationRequest;
use RuntimeException;

class RunOperation
{
    /**
     * @throws  ValidationException
     */
    public function doOperation(): ResultOperation
    {
        //создаем сотрудников для наглядности
        $this->prepareEmployees();

        //берем данные и валидируем, конечно тут частичная валидация
        //и полная проверка должна быть до использования данных
        $data = RequestMapper::fromPost();
        //валидируем
        ValidationRequest::validate($data);

        /** @var Employee $reseller */
        $reseller = Employee::getById(type: EmployeeTypeEnum::RESELLER, id: $data->resellerId);
        /** @var Employee $client */
        $client = Employee::getById(type: EmployeeTypeEnum::CONTRACTOR, id: $data->clientId);
        /** @var Employee $creator */
        $creator = Employee::getById(type: EmployeeTypeEnum::CREATOR, id: $data->creatorId);
        /** @var Employee $expert */
        $expert = Employee::getById(type: EmployeeTypeEnum::EXPERT, id: $data->expertId);


        $differences = '';
        if ($data->notificationType === NotificationTypeEnum::NEW) {
            $differences = __('NewPositionAdded');
        } elseif ($data->notificationType === NotificationTypeEnum::CHANGE) {
            $params = null;
            if ($data->differences) {
                $params = [
                    'FROM' => $data->differences->from->name,
                    'TO'   => $data->differences->to->name,
                ];
            }

            $differences = __('PositionStatusHasChanged', $params);
        }

        //Если $differences пустой, значит появился новый тип в NotificationTypeEnum
        if (empty($differences)) {
            throw new RuntimeException("Template Data (differences) is empty!", 500);
        }

        $templateData = [
            'COMPLAINT_ID'       => $data->complaintId,
            'COMPLAINT_NUMBER'   => $data->complaintNumber,
            'CREATOR_ID'         => $creator->getId(),
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $expert->getId(),
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $client->getId(),
            'CLIENT_NAME'        => $client->getFullName(),
            'CONSUMPTION_ID'     => $data->consumptionId,
            'CONSUMPTION_NUMBER' => $data->consumptionNumber,
            'AGREEMENT_NUMBER'   => $data->agreementNumber,
            'DATE'               => $data->date,
            'DIFFERENCES'        => $differences,
        ];


        $resultOperation = new ResultOperation();

        $emailFrom = $reseller->getEmail();
        // Получаем email сотрудников из настроек
        $emails = Emails::getEmailsByPermit();

        if ($emailFrom && ! empty($emails)) {
            $buff = [];
            $subject = __('complaintEmployeeEmailSubject', $templateData);
            $message = __('complaintEmployeeEmailBody', $templateData);
            foreach ($emails as $email) {
                $buff[] = [
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $email,
                    'subject'   => $subject,
                    'message'   => $message,
                ];
            }
            MessagesClient::sendMessage($buff, $reseller->getId(),
                NotificationsEvents::NEW_RETURN_STATUS);
            $resultOperation->notificationEmployeeByEmailSent();
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($data->notificationType === NotificationTypeEnum::CHANGE && $data->differences->to) {
            if ($emailFrom && $client->getEmail()) {
                MessagesClient::sendMessage([
                    [
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $client->getEmail(),
                        'subject'   => __('complaintClientEmailSubject', $templateData),
                        'message'   => __('complaintClientEmailBody', $templateData),
                    ],
                ], $reseller->getId(),
                    $client->getId(),
                    NotificationsEvents::CHANGE_RETURN_STATUS,
                    $data->differences->to);
                $resultOperation->notificationClientByEmailSent();
            }

            if ($client->getPhone()) {
                $res = NotificationManager::send($reseller->getId(), $client->getId(),
                    NotificationsEvents::CHANGE_RETURN_STATUS, $data->differences->to,
                    $templateData, $error);
                if ($res) {
                    $resultOperation->notificationClientBySmsSent();
                }
                if ( ! empty($error)) {
                    $resultOperation->setNotificationClientBySmsMess($error);
                }
            }
        }

        return $resultOperation;
    }

    //создаем пользователей
    private function prepareEmployees(): void
    {
        Employee::setEmployee(
            new Employee(id: 1,
                type: EmployeeTypeEnum::RESELLER,
                name: 'Name',
                surname: 'SurName',
                email: 'example@gmail.com',
                phone: "+7777777777")
        );

        Employee::setEmployee(
            new Employee(id: 1,
                type: EmployeeTypeEnum::CONTRACTOR,
                name: 'Name',
                surname: 'SurName',
                email: 'example@gmail.com',
                phone: "+7777777777")
        );

        Employee::setEmployee(
            new Employee(id: 1,
                type: EmployeeTypeEnum::CREATOR,
                name: 'Name',
                surname: 'SurName',
                email: 'example@gmail.com',
                phone: "+7777777777")
        );

        Employee::setEmployee(
            new Employee(id: 1,
                type: EmployeeTypeEnum::EXPERT,
                name: 'Name',
                surname: 'SurName',
                email: 'example@gmail.com',
                phone: "+7777777777")
        );
    }
}
