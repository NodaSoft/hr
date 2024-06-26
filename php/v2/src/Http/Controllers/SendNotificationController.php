<?php

namespace Nodasoft\Testapp\Http\Controllers;

use Nodasoft\Testapp\DTO\MessageDifferenceDto;
use Nodasoft\Testapp\DTO\SendNotificationDTO;
use Nodasoft\Testapp\Enums\NotificationType;
use Nodasoft\Testapp\Http\Responses\JsonResponse;
use Nodasoft\Testapp\Services\SendNotification\Base\ReferencesOperation;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;

final readonly class SendNotificationController
{
    public function __construct(private ReferencesOperation $referencesOperation)
    {
    }

    public function send(ServerRequestInterface $request): ResponseInterface
    {
        $data = $request->getParsedBody();

        // requirement:--
        // if differences isset in request then
        // is required to pass also from and to (status id) field
        $difference = isset($data['differences']) ? new MessageDifferenceDto(
            $data['differences']['from'],
            $data['differences']['to'],
        ) : null;

        $this->referencesOperation->doOperation(new SendNotificationDTO(
            $data['reseller_id'],
            NotificationType::tryFrom($data['notification_type']),
            $data['client_id'],
            $data['creator_id'],
            $data['expert_id'],
            $data['complaint_id'],
            $data['complaint_number'],
            $data['consumption_id'],
            $data['consumption_number'],
            $data['agreement_number'],
            $data['date'],
            $difference
        ));

        return JsonResponse::json([], 'handled successfully!');
    }
}