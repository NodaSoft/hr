<?php

declare(strict_types=1);

namespace ResultOperation\EventListener;

use ResultOperation\Enum\NotificationEvent;
use ResultOperation\Event\AbstractStatusEvent;
use ResultOperation\Event\ChangeStatusEvent;
use ResultOperation\Event\NewStatusEvent;

use ResultOperation\Service\MessagesClient;
use ResultOperation\Service\NotificationManager;

use function NW\WebService\References\Operations\Notification\getEmailsByPermit;
use function NW\WebService\References\Operations\Notification\getResellerEmailFrom;

class NotificationEventListener
{
    public function __construct(
        private readonly MessagesClient $messagesClient,
        private readonly NotificationManager $notificationManager
    ) {
    }

    /**
     * @var array
     */
    private const SUBSCRIBED_EVENT = [
        NotificationEvent::NEW->value => 'handleNewStatus',
        NotificationEvent::CHANGE->value => 'handleChangeStatus',
    ];

    /**
     * @return array
     */
    public function getSubscribedEvents(): array
    {
        return self::SUBSCRIBED_EVENT;
    }

    /**
     * @param NewStatusEvent $event
     * @return NewStatusEvent
     */
    public function handleNewStatus(NewStatusEvent $event): NewStatusEvent
    {
        /** @var NewStatusEvent $event */
        $event = $this->sendEmployersEmails($event, NotificationEvent::NEW);

        return $event;
    }

    /**
     * @param ChangeStatusEvent $event
     * @return ChangeStatusEvent
     */
    public function handleChangeStatus(ChangeStatusEvent $event): ChangeStatusEvent
    {
        /** @var ChangeStatusEvent $event */
        $event = $this->sendEmployersEmails($event, NotificationEvent::CHANGE);

        $client = $event->getClient();
        $notificationTemplate = $event->getTemplate();
        $templateData = $notificationTemplate->toArray();

        $resellerId = $notificationTemplate->getResellerId();
        $emailFrom = getResellerEmailFrom($resellerId);

        if (!empty($emailFrom) && !empty($client->getEmail())) {
            $this->messagesClient->sendMessage(
                [
                    [
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $client->getEmail(),
                        'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ],
                $resellerId,
                $client->getId(),
                NotificationEvent::CHANGE->value,
                (int) $templateData['differences']['to']
            );
            $event->setClientByEmail();
        }

        if (!empty($client->hasMobile())) {
            $this->notificationManager->send(
                $resellerId,
                $client->getId(),
                NotificationEvent::CHANGE->value,
                (int) $templateData['differences']['to'],
                $templateData,
                $error
            );

            if (!empty($error)) {
                return $event->setError($error);
            }

            $event->setClientBySms();
        }

        return $event;
    }

    /**
     * @param AbstractStatusEvent $event
     * @param NotificationEvent $eventType
     * @return AbstractStatusEvent
     */
    public function sendEmployersEmails(AbstractStatusEvent $event, NotificationEvent $eventType): AbstractStatusEvent
    {
        $notificationTemplate = $event->getTemplate();
        $templateData = $notificationTemplate->toArray();

        $resellerId = $notificationTemplate->getResellerId();
        $emailFrom = getResellerEmailFrom($resellerId);

        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                $this->messagesClient->sendMessage(
                    [
                        [
                            'emailFrom' => $emailFrom,
                            'emailTo'   => $email,
                            'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                            'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                        ],
                    ],
                    $resellerId,
                    $eventType->value
                );
                $event->setEmployeeByEmail();
            }
        }

        return $event;
    }
}
