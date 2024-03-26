<?php

namespace Src\Operation\Infrastructure\Adapters;

use Src\Contractor\Infrastructure\API\ContractorApi;
use Src\Operation\Application\DataTransferObject\ContractorData;
use src\Operation\Application\Exceptions\ContractorNotFoundException;

readonly class ContractorAdapter
{
    private ContractorApi $contractorApi;

    public function __construct()
    {
        $this->contractorApi = new ContractorApi();
    }

    /**
     * @throws ContractorNotFoundException
     */
    public function getById(int $contractorId): ContractorData
    {
        /**
         * Здесь могла бы быть обработка ошибок, таких как:
         * - входные данные для запроса к сервису-провайдеру некорректны или не полны.
         * - не удается установить соединение с сервисом-провайдером (например, из-за проблем с сетью).
         * - полученный ответ от сервиса-провайдера не соответствует ожидаемому формату или содержит ошибки.
         *
         * Данный сервис-провайдер возвращает ошибку, которую необходимо обработать на стороне клиента
         * (например, переформулировать сообщение об ошибке для конечного пользователя).
         */

        $contractor = $this->contractorApi->getById($contractorId);

        if ($contractor == null) {
            throw new ContractorNotFoundException('Contractor not found', 400);
        }

        return ContractorData::fromArray($this->contractorApi->getById($contractorId));
    }

    public function getEmailsByPermit($resellerId, $event): array
    {
        return $this->contractorApi->getEmailsByPermit($resellerId, $event);
    }

}