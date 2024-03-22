<?php

namespace Src\Employee\Infrastructure\API;

use EmployeeData;
use Src\Employee\Infrastructure\Repository\EmployeeRepository;

class EmployeeApi
{
    private EmployeeRepository $repository;

    public function __construct()
    {
        $this->repository = new EmployeeRepository();

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

        $employee = $this->repository->getById($id);
        return EmployeeData::fromEntity($employee)->toArray();
    }

}