<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

use Exception;
use NW\WebService\References\Operations\Notification\Client\MessagesClient;
use NW\WebService\References\Operations\Notification\Client\NotificationManager;
use NW\WebService\References\Operations\Notification\Domain\Client;
use NW\WebService\References\Operations\Notification\Domain\Employee;
use NW\WebService\References\Operations\Notification\Domain\ReferencesOperation;
use NW\WebService\References\Operations\Notification\Event\ChangeReturnStatusEvent;
use NW\WebService\References\Operations\Notification\Struct\Differences;
use NW\WebService\References\Operations\Notification\Struct\Email;
use NW\WebService\References\Operations\Notification\Struct\Result;
use NW\WebService\References\Operations\Notification\Struct\SmsNotification;
use NW\WebService\References\Operations\Notification\Struct\Template;
use Symfony\Component\HttpFoundation\Exception\BadRequestException;
use Symfony\Component\HttpFoundation\Response;

/**
 * @method newPositionAdded(int $sellerId, ?Differences $differences = null)
 * @method positionStatusHasChanged(Differences $differences, int $sellerId)
 * @method complaintClientEmailSubject(Template $template, int $sellerId)
 * @method complaintClientEmailBody(Template $template, int $sellerId)
 */
class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws Exception
     */
    public function doOperation(): Result
    {
        $data = $this->getRequest('data');

        $client = Client::find($data->clientId);

        if ($client->seller->id !== $data->resellerId) {
            throw new Exception('Client not found!', Response::HTTP_BAD_REQUEST);
        }

        $differences = match ($data->notificationType) {
            self::TYPE_NEW => $this->newPositionAdded($client->seller->id),
            /** @phpstan-ignore-next-line */
            self::TYPE_CHANGE => $this->positionStatusHasChanged($data->differences, $client->seller->id),
            default           => null
        };

        $creator = Employee::find($data->creatorId);

        $expert = Employee::find($data->expertId);

        $template = new Template(
            complaintId: $data->complaintId,
            complaintNumber: $data->complaintNumber,
            creatorId: $creator->id,
            creatorName: $creator->getFullName(),
            expertId:  $expert->id,
            expertName: $expert->getFullName(),
            clientId:  $client->id,
            clientName: $client->getFullName(),
            consumptionId:  $data->consumptionId,
            consumptionNumber:  $data->consumptionNumber,
            agreementNumber: $data->agreementNumber,
            date: $data->date,
            differences: $differences ?? null
        );

        $emailFrom = $client->seller->email;

        // Получаем email сотрудников из настроек
        $emails = $client->seller->getEmails();

        $result = new Result();

        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage(
                    email: new Email(
                        $emailFrom,
                        $email,
                        $this->complaintClientEmailSubject($template, $client->seller->id),
                        $this->complaintClientEmailBody($template, $client->seller->id)
                    ),
                    sellerId: $client->seller->id,
                    event: new ChangeReturnStatusEvent()
                );

                $result->notifiedEmployeeByEmail();
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($data->notificationType === self::TYPE_CHANGE && $data->differences !== null) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage(
                    email: new Email(
                        $emailFrom,
                        $client->email,
                        $this->complaintClientEmailSubject($template, $client->seller->id),
                        $this->complaintClientEmailBody($template, $client->seller->id)
                    ),
                    sellerId: $client->seller->id,
                    event: new ChangeReturnStatusEvent(),
                    clientId: $client->id,
                    status: $data->differences->to
                );

                $result->notifiedClientByEmail();
            }

            if (!empty($client->mobile)) {

                try {
                    NotificationManager::send(
                        $client->seller->id,
                        $client->id,
                        new ChangeReturnStatusEvent(),
                        $data->differences->to,
                        $template
                    );
                } catch (BadRequestException $exception) {
                    return $result->setSmsNotification(new SmsNotification(
                        isSent: false,
                        message: $exception->getMessage()
                    ));
                }

                return $result->setSmsNotification(new SmsNotification(
                    isSent: true,
                    message: 'Success'
                ));
            }
        }

        return $result;
    }
}
