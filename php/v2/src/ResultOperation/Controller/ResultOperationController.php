<?php

declare(strict_types=1);

namespace ResultOperation\Controller;

use Exception;
use pseudovendor\BaseController;
use pseudovendor\ErrorResponse;
use pseudovendor\EventDispatcherInterface;
use pseudovendor\Request;
use pseudovendor\Response;
use pseudovendor\SuccessResponse;
use ResultOperation\DTO\NotificationTemplate;
use ResultOperation\Entity\Contractor;
use ResultOperation\Entity\Employee;
use ResultOperation\Entity\Seller;
use ResultOperation\Enum\ContractorType;
use ResultOperation\Enum\NotificationEvent;
use ResultOperation\Enum\Status;
use ResultOperation\Event\ChangeStatusEvent;
use ResultOperation\Event\NewStatusEvent;
use ResultOperation\Exception\AbstractEntityNotFoundByIdException;
use ResultOperation\Exception\ClientNotFoundException;
use ResultOperation\Exception\CreatorNotFoundException;
use ResultOperation\Exception\ExpertNotFoundException;
use ResultOperation\Exception\InvalidClientException;
use ResultOperation\Exception\ResellerNotFountException;
use ResultOperation\Exception\SellerNotFountException;
use ResultOperation\Repository\ContractorRepository;
use ResultOperation\Repository\EmployeeRepository;
use ResultOperation\Repository\SellerRepository;

class ReturnOperationController extends BaseController
{
    /**
     * Фронтовые коды для notificationType
     */
    private const FRONT_NEW_STATUS_CODE = 1;
    private const FRONT_CHANGE_STATUS_CODE = 2;

    public function __construct(
        private readonly EventDispatcherInterface $eventDispatcher,
        private readonly ContractorRepository $contractorRepository,
        private readonly EmployeeRepository $employeeRepository,
        private readonly SellerRepository $sellerRepository
    ) {
    }

    /**
     * @param Request $request
     * @return Response
     * @throws Exception
     */
    public function index(Request $request): Response
    {
        $data = $request->getBody()['data'] ?? null;
        if (!is_array($data)) {
            return new ErrorResponse(/*bla bla*/);
        }

        $errors = $this->validate($data);
        if ($errors !== null) {
            return $errors;
        }

        try {
            $clientId = (int) ($data['clientId'] ?? 0);
            /** @var Contractor $client */
            $client = $this->contractorRepository->get($clientId);
            if ($client === null) {
                throw new ClientNotFoundException($clientId);
            }

            $notificationType = (int) ($data['notificationType'] ?? null);
            switch ($notificationType) {
                case self::FRONT_NEW_STATUS_CODE:
                    $template = $this->buildTemplate($client, $data, $notificationType);

                    /** @var NewStatusEvent $result */
                    $result = $this->eventDispatcher->dispatch(
                        NotificationEvent::NEW->value,
                        new NewStatusEvent($client, $template)
                    );
                    break;
                case self::FRONT_CHANGE_STATUS_CODE:
                    $template = $this->buildTemplate($client, $data, $notificationType);

                    /** @var ChangeStatusEvent $result */
                    $result = $this->eventDispatcher->dispatch(
                        NotificationEvent::CHANGE->value,
                        new ChangeStatusEvent($client, $template)
                    );
                    break;
            }
        } catch (AbstractEntityNotFoundByIdException $exception) {
            return new ErrorResponse(/*bla bla*/);
        }

        if (!isset($result) || $result->getError()) {
            return new ErrorResponse(/*bla bla*/);
        }

        /**
         * ...манипуляции с @var $result
         */

        return new SuccessResponse();
    }

    /**
     * Провалидирует входные данные и в случае ошибки вернет ErrorResponse
     *
     * Если речь идет не о частном запросе, а, например, об апишке,
     * то естественно все это чудо надо выносить в валидатор
     *
     * @param array $data
     * @return ?ErrorResponse
     */
    private function validate(array $data): ?ErrorResponse
    {
        return null; // stub
    }

    /**
     * Раз уж у нас один шаблон и на клиентское уведомление, и на почтовое,
     * то и формироваться он будет в одном месте
     *
     * @param Contractor $client
     * @param array $data
     * @param int $notificationType
     *
     * @return NotificationTemplate
     *
     * @throws ClientNotFoundException
     * @throws CreatorNotFoundException
     * @throws ExpertNotFoundException
     * @throws InvalidClientException
     * @throws ResellerNotFountException
     * @throws SellerNotFountException
     */
    private function buildTemplate(Contractor $client, array $data, int $notificationType): NotificationTemplate
    {
        $resellerId = (int) ($data['resellerId'] ?? 0);
        /** @var Seller $reseller */
        $reseller = $this->sellerRepository->get($resellerId);
        if ($reseller === null) {
            throw new ResellerNotFountException($resellerId);
        }

        if ($client->getType() !== ContractorType::CUSTOMER) {
            throw new InvalidClientException(
                sprintf(
                    'Invalid client type: %s',
                    $client->getType()->name
                )
            );
        }

        $sellerId = (int) $client->getSellerId();
        if ($sellerId === 0) {
            throw new InvalidClientException('Client has no linked seller');
        }
        if ($this->sellerRepository->get($sellerId) === null) {
            throw new SellerNotFountException($sellerId);
        }

        $creatorId = (int) ($data['creatorId'] ?? 0);
        /** @var Employee $creator */
        $creator = $this->employeeRepository->get($creatorId);
        if ($creator === null) {
            throw new CreatorNotFoundException($creatorId);
        }

        $expertId = (int) ($data['expertId'] ?? 0);
        /** @var Employee $expert */
        $expert = $this->employeeRepository->get($expertId);
        if ($expert === null) {
            throw new ExpertNotFoundException($expertId);
        }

        /**
         * Что делает метод @see __() для меня осталось загадкой,
         * поэтому пробросил {@var $notificationType} и этот участок кода не трогал
         */
        $differences = '';
        if ($notificationType === self::FRONT_NEW_STATUS_CODE) {
            $differences = __(
                'NewPositionAdded',
                null,
                $resellerId
            );
        } elseif ($notificationType === self::FRONT_CHANGE_STATUS_CODE && !empty($data['differences'])) {
            $differences = __(
                'PositionStatusHasChanged',
                [
                    'FROM' => Status::from((int) $data['differences']['from']),
                    'TO'   => Status::from((int) $data['differences']['to']),
                ],
                $resellerId
            );
        }

        return (new NotificationTemplate())
            ->setComplaintId((int) $data['complaintId'])
            ->setResellerId($resellerId)
            //...
        ;
    }
}
