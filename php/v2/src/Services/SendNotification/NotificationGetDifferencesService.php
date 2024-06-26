<?php

namespace Nodasoft\Testapp\Services\SendNotification;

use Exception;
use Nodasoft\Testapp\DTO\GetNotificationDifferenceDTO;
use Nodasoft\Testapp\Repositories\Status\StatusRepositoryInterface;
use Nodasoft\Testapp\Services\SendNotification\Base\GetDifferencesInterface;

readonly class NotificationGetDifferencesService implements GetDifferencesInterface
{
    public function __construct(private StatusRepositoryInterface $statusRepository)
    {
    }

    /**
     * @throws Exception
     */
    public function getDifference(GetNotificationDifferenceDTO $differenceDTO)
    {
        if ($differenceDTO->notificationType->isNew()) {
            return __(
                'NewPositionAdded',
                null,
                $differenceDTO->resellerId
            );
        }

        if ($differenceDTO->messageDifference) {
            return __(
                'PositionStatusHasChanged',
                [
                    'FROM' => $this->statusRepository->getById($differenceDTO->messageDifference->from)->getName(),
                    'TO' => $this->statusRepository->getById($differenceDTO->messageDifference->to)->getName(),
                ],
                $differenceDTO->resellerId
            );
        }

        return '';
    }
}