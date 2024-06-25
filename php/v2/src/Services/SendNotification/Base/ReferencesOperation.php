<?php

namespace Nodasoft\Testapp\Services\SendNotification\Base;


use Nodasoft\Testapp\DTO\SendNotificationDTO;

interface  ReferencesOperation
{
    public function doOperation(SendNotificationDTO $dto): array;
}