<?php

namespace Nodasoft\Testapp\Services\SendNotification\Base;

use Nodasoft\Testapp\DTO\GetNotificationDifferenceDTO;

interface GetDifferencesInterface
{
    public function getDifference(GetNotificationDifferenceDTO $differenceDTO);
}