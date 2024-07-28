<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;
use NW\WebService\References\Operations\Notification\Mailer\NotificationManager;
use NW\WebService\References\Operations\Notification\Notification\Result;
use NW\WebService\References\Operations\Notification\Validation\Rule\ClientValidator;
use NW\WebService\References\Operations\Notification\Validation\Rule\EmployeeValidator;
use NW\WebService\References\Operations\Notification\Validation\Rule\RequestValidator;
use NW\WebService\References\Operations\Notification\Validation\Rule\SellerValidator;
use NW\WebService\References\Operations\Notification\Validation\ValidationBuilder;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    protected NotificationManager $notificationManager;

    public function __construct(NotificationManager $notificationManager)
    {
        $this->notificationManager = $notificationManager;
    }


    /**
     * @return array
     * @throws Exception
     */
    public function doOperation(): array
    {
        $data = (array)$this->getRequest('data');
        $result = (new Result())->initialize();

        $validationBuilder = new ValidationBuilder();
        $validationBuilder
            ->add(new RequestValidator())
            ->add(new SellerValidator())
            ->add(new EmployeeValidator())
            ->add(new ClientValidator());

        if (!$validationBuilder->validate($data, $result)) {
            return $result;
        }

        $this->notificationManager->sendEmployeeNotification($data, $result);

        if ($data['notificationType'] === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            $this->notificationManager->sendClientNotification($data, $result);
        }

        return $result;
    }
}
