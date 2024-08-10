<?php

namespace app\Http\Controllers\References\Operations\Notification;

use app\Domain\Notification\Actions\NotificationAction;
use app\Domain\Notification\Exceptions\ActionException;
use App\Http\Controllers\Notification\DataTransferObjectError;
use app\Http\Controllers\References\Operations\OperationController;
use app\Domain\Notification\DTO\NotificationData;

class NotificationController extends OperationController
{
    /**
     * @throws \Exception
     * @return array
     */
    public function doOperation(): array
    {
        $result = [];
        try {
            $data = NotificationData::fromNotificationRequest((array)$this->getRequest('data'));
            $result = (new NotificationAction())->notify($data);
        } catch (DataTransferObjectError $e) {
            throw new \Exception($e->getMessage(), 500);
        } catch (ActionException $ae) {
            throw new ActionException($ae->getMessage(), $ae->getCode());
        }
        return $result;
    }
}