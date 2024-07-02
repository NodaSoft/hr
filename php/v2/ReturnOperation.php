<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\DTO\Notification\OperationResultDTO;
use Exception;
use NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException;
use NW\WebService\References\Operations\Notification\Helpers\NotificationEvents;
use NW\WebService\References\Operations\Notification\Helpers\NotificationsManagement;
use NW\WebService\References\Operations\Notification\Helpers\ReferencesOperation;
use NW\WebService\References\Operations\Notification\Helpers\Status;
use NW\WebService\References\Operations\Notification\Helpers\Support;
use NW\WebService\References\Operations\Notification\Interfaces\LoggerInterface;
use Repositories\ContractorRepository;

/**
 * Class TsReturnOperation
 * Handles the TS Goods Return operation including Email and SMS notifications.
 */
class TsReturnOperation extends ReferencesOperation
{
    /** @var string */
    public const EVENT_TYPE = 'tsGoodsReturn';

    /** @var \Repositories\ContractorRepository */
    private readonly ContractorRepository $contractorRepository;

    /** @var \NW\WebService\References\Operations\Notification\Helpers\NotificationsManagement */
    private readonly NotificationsManagement $notificationsManagement;

    /**
     * Mock interface without implementation. Use Psr Logger in Production
     *
     * @var \NW\WebService\References\Operations\Notification\Interfaces\LoggerInterface
     */
    private readonly LoggerInterface $logger;

    public function __construct()
    {
        $this->contractorRepository = new ContractorRepository();
        $this->notificationsManagement = new NotificationsManagement();
    }

    /** @var int */
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @return \NW\WebService\References\Operations\Notification\DTO\Notification\OperationResultDTO
     * @throws \NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException
     */
    public function doOperation(): OperationResultDTO
    {
        $data = $this->getRequest('data');
        $reseller = null;
        $client = null;
        $creator = null;
        $expert = null;

        try {
            /**
             * @todo If all these Entities are the same DB Resources - can be combined into Collection instead.
             */
            $reseller = $this->contractorRepository->getById((int) $data['resellerId']);
            $client = $this->contractorRepository->getById((int) $data['clientId']);
            $creator = $this->contractorRepository->getById((int) $data['creatorId']);
            $expert = $this->contractorRepository->getById((int) $data['expertId']);
        } catch (Exception $e) {
            $this->logger->log(
                $e->getMessage(),
                [
                    'data' => $data,
                ]
            );
            $entityName = match (null) {
                $reseller => 'Seller',
                $client => 'Client',
                $creator => 'Creator',
                $expert => 'Expert',
                default => null,
            };
            if ($entityName) {
                throw new InvalidArgumentsException(
                    Support::__(':entityName not Found!', [':entityName' => $entityName]),
                    400
                );
            }
        }

        $result = new OperationResultDTO();
        $templateData = [
            'EVENT_STATUS'       => NotificationEvents::CHANGE_RETURN_STATUS,
            'COMPLAINT_ID'       => (int) $data['complaintId'],
            'COMPLAINT_NUMBER'   => (string) $data['complaintNumber'],
            'CREATOR_ID'         => (int) $data['creatorId'],
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => (int)$data['expertId'],
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => (int)$data['clientId'],
            'CLIENT_NAME'        => $client->getFullName(),
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $this->getDifferences($data),
        ];

        try {
            $sendResult = $this->notificationsManagement->emailSend(
                $reseller->getEmail(),
                $reseller->getEmailsByPermit(self::EVENT_TYPE),
                $templateData,
                'complaintEmployeeEmailSubject',
                'complaintEmployeeEmailBody',
            );
            $result->setIsEmployeeNotifiedByEmail($sendResult);
        } catch (Exception $e) {
            $this->logger->log($e->getMessage(), ['data' => $templateData]);
            $result->setEmployeeEmailNotificationMessage($e->getMessage());
        }

        $notificationType = (int) $data['notificationType'];

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            try {
                $sendResult = $this->notificationsManagement->emailSend(
                    $reseller->getEmail(),
                    $client->getEmail(),
                    $templateData,
                    'complaintClientEmailSubject',
                    'complaintClientEmailBody'
                );
                $result->setIsClientNotifiedByEmail($sendResult);
            } catch (Exception $e) {
                $this->logger->log($e->getMessage(), ['data' => $templateData]);
                $result->setClientEmailNotificationMessage($e->getMessage());
            }

            try {
                $sendResult = $this->notificationsManagement->smsSend(
                    $client->getMobile(),
                    $reseller->getId(),
                    $client->getId(),
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $templateData,
                );
                $result->setIsClientNotifiedBySms($sendResult);
            } catch (Exception $e) {
                $this->logger->log($e->getMessage(), ['data' => $templateData]);
                $result->setClientSmsNotificationMessage($e->getMessage());
            }
        }

        return $result;
    }

    /**
     * Generates a string representation of the differences based on the provided data.
     *
     * @param array $data
     * @return string
     * @throws \NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException
     */
    private function getDifferences(array $data): string
    {
        $notificationType = (int) $data['notificationType'];
        if (!$notificationType) {
            throw new InvalidArgumentsException('Empty notificationType', 400);
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = Support::__('NewPositionAdded');
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = Support::__('PositionStatusHasChanged', [
                'FROM' => Status::getName((int) ($data['differences']['from'] ?? 0)),
                'TO'   => Status::getName((int) ($data['differences']['to'] ?? 0)),
            ]);
        }

        return $differences;
    }
}
