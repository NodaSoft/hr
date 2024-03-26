<?php

namespace Src\Contractor\Infrastructure\API;

use Src\Contractor\Application\DataTransferObject\ContractorData;
use Src\Contractor\Infrastructure\Repository\NotificationRepository;

class ContractorApi
{
    private NotificationRepository $repository;

    public function __construct()
    {
        $this->repository = new NotificationRepository();

    }

    public function getById(int $id): array
    {
        /**
         * - так как ContractorApi выступает в качестве провайдера данных, то здесь нужна валидация входных данных
         * - ошибки доступа к ресурсам, например к базе данных, должны обрабатываться здесь
         * - невозможность выполнения запрошенной операции из-за внутренних ограничений или ошибок.
         * - проблемы с зависимостями или внешними сервисами, от которых зависит провайдер
         *
         * Данный сервис должен возвщать информативные сообщения об ошибках(в рамках безопасности, чтобы не раскрывать
         * внутренние детали, позволяя клиентскому сервису адекватно реагировать на эти ошибки)
         */

        $contractor = $this->repository->getById($id);
        return ContractorData::fromEntity($contractor)->toArray();
    }

    public function getEmailsByPermit($resellerId, $event): array
    {
        // fakes the method
        return ['someemeil@example.com', 'someemeil2@example.com'];
    }

}